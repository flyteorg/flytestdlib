package storage

import (
	"bytes"
	"context"
	"io"
	"runtime/debug"
	"time"

	"github.com/flyteorg/flytestdlib/errors"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/flyteorg/flytestdlib/promutils"

	"github.com/coocood/freecache"
	"github.com/flyteorg/flytestdlib/ioutils"
	"github.com/flyteorg/flytestdlib/logger"
)

const neverExpire = 0

// TODO Freecache has bunch of metrics it calculates. Lets write a prom collector to publish these metrics
type cacheMetrics struct {
	CacheHit        prometheus.Counter
	CacheMiss       prometheus.Counter
	CacheWriteError prometheus.Counter
	FetchLatency    promutils.StopWatch
}

type cachedRawStore struct {
	RawStore
	cache   *freecache.Cache
	scope   promutils.Scope
	metrics *cacheMetrics
}

// Gets metadata about the reference. This should generally be a light weight operation.
func (s *cachedRawStore) Head(ctx context.Context, reference DataReference) (Metadata, error) {
	key := []byte(reference)
	if oRaw, err := s.cache.Get(key); err == nil {
		s.metrics.CacheHit.Inc()
		// Found, Cache hit
		size := int64(len(oRaw))
		// return size in metadata
		return StowMetadata{exists: true, size: size}, nil
	}
	s.metrics.CacheMiss.Inc()
	return s.RawStore.Head(ctx, reference)
}

// Retrieves a byte array from the Blob store or an error
func (s *cachedRawStore) ReadRaw(ctx context.Context, reference DataReference) (io.ReadCloser, error) {
	key := []byte(reference)
	if oRaw, err := s.cache.Get(key); err == nil {
		// Found, Cache hit
		s.metrics.CacheHit.Inc()
		return ioutils.NewBytesReadCloser(oRaw), nil
	}
	s.metrics.CacheMiss.Inc()
	reader, err := s.RawStore.ReadRaw(ctx, reference)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = reader.Close()
		if err != nil {
			logger.Warnf(ctx, "Failed to close reader [%v]. Error: %v", reference, err)
		}
	}()

	b, err := ioutils.ReadAll(reader, s.metrics.FetchLatency.Start())
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(key, b, 0)
	if err != nil {
		logger.Debugf(ctx, "Failed to Cache the metadata")
		err = errors.Wrapf(ErrFailedToWriteCache, err, "Failed to Cache the metadata")
	}

	return ioutils.NewBytesReadCloser(b), err
}

// Stores a raw byte array.
func (s *cachedRawStore) WriteRaw(ctx context.Context, reference DataReference, size int64, opts Options, raw io.Reader) error {
	var buf bytes.Buffer
	teeReader := io.TeeReader(raw, &buf)
	err := s.RawStore.WriteRaw(ctx, reference, size, opts, teeReader)
	if err != nil {
		return err
	}

	err = s.cache.Set([]byte(reference), buf.Bytes(), neverExpire)
	if err != nil {
		s.metrics.CacheWriteError.Inc()
		err = errors.Wrapf(ErrFailedToWriteCache, err, "Failed to Cache the metadata")
	}

	return err
}

// Creates a CachedStore if Caching is enabled, otherwise returns a RawStore
func newCachedRawStore(cfg *Config, store RawStore, scope promutils.Scope) RawStore {
	if cfg.Cache.MaxSizeMegabytes > 0 {
		c := &cachedRawStore{
			RawStore: store,
			cache:    freecache.NewCache(cfg.Cache.MaxSizeMegabytes * 1024 * 1024),
			scope:    scope,
			metrics: &cacheMetrics{
				FetchLatency:    scope.MustNewStopWatch("remote_fetch", "Total Time to read from remote metastore", time.Millisecond),
				CacheHit:        scope.MustNewCounter("cache_hit", "Number of times metadata was found in cache"),
				CacheMiss:       scope.MustNewCounter("cache_miss", "Number of times metadata was not found in cache and remote fetch was required"),
				CacheWriteError: scope.MustNewCounter("cache_write_err", "Failed to write to cache"),
			},
		}
		if cfg.Cache.TargetGCPercent > 0 {
			debug.SetGCPercent(cfg.Cache.TargetGCPercent)
		}
		return c
	}
	return store
}
