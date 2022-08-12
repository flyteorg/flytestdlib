package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/promutils/labeled"
)

type dataStoreCreateFn func(cfg *Config, metrics *DataStoreMetrics) (RawStore, error)

var stores = map[string]dataStoreCreateFn{
	TypeMemory: NewInMemoryRawStore,
	TypeLocal:  newStowRawStore,
	TypeMinio:  newStowRawStore,
	TypeS3:     newStowRawStore,
	TypeStow:   newStowRawStore,
}

type proxyTransport struct {
	http.RoundTripper
	defaultHeaders map[string][]string
}

func (p proxyTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	applyDefaultHeaders(r, p.defaultHeaders)
	return p.RoundTripper.RoundTrip(r)
}

func applyDefaultHeaders(r *http.Request, headers map[string][]string) {
	if r.Header == nil {
		r.Header = http.Header{}
	}

	for key, values := range headers {
		for _, val := range values {
			r.Header.Add(key, val)
		}
	}
}

func createHTTPClient(cfg HTTPClientConfig) *http.Client {
	c := &http.Client{
		Timeout: cfg.Timeout.Duration,
	}

	if len(cfg.Headers) > 0 {
		c.Transport = &proxyTransport{
			RoundTripper:   http.DefaultTransport,
			defaultHeaders: cfg.Headers,
		}
	}

	return c
}

type DataStoreMetrics struct {
	cacheMetrics *cacheMetrics
	protoMetrics *protoMetrics
	copyMetrics  *copyMetrics
	stowMetrics  *stowMetrics
}

// NewDataStoreMetrics initialises all metrics required for DataStore
func NewDataStoreMetrics(scope promutils.Scope) *DataStoreMetrics {
	failureTypeOption := labeled.AdditionalLabelsOption{Labels: []string{FailureTypeLabel.String()}}
	return &DataStoreMetrics{
		cacheMetrics: &cacheMetrics{
			FetchLatency:    scope.MustNewStopWatch("remote_fetch", "Total Time to read from remote metastore", time.Millisecond),
			CacheHit:        scope.MustNewCounter("cache_hit", "Number of times metadata was found in cache"),
			CacheMiss:       scope.MustNewCounter("cache_miss", "Number of times metadata was not found in cache and remote fetch was required"),
			CacheWriteError: scope.MustNewCounter("cache_write_err", "Failed to write to cache"),
		},
		protoMetrics: &protoMetrics{
			FetchLatency:                 scope.MustNewStopWatch("proto_fetch", "Time to read data before unmarshalling", time.Millisecond),
			MarshalTime:                  scope.MustNewStopWatch("marshal", "Time incurred in marshalling data before writing", time.Millisecond),
			UnmarshalTime:                scope.MustNewStopWatch("unmarshal", "Time incurred in unmarshalling received data", time.Millisecond),
			MarshalFailure:               scope.MustNewCounter("marshal_failure", "Failures when marshalling"),
			UnmarshalFailure:             scope.MustNewCounter("unmarshal_failure", "Failures when unmarshalling"),
			WriteFailureUnrelatedToCache: scope.MustNewCounter("write_failure_unrelated_to_cache", "Raw store write failures that are not caused by ErrFailedToWriteCache"),
			ReadFailureUnrelatedToCache:  scope.MustNewCounter("read_failure_unrelated_to_cache", "Raw store read failures that are not caused by ErrFailedToWriteCache"),
		},
		copyMetrics: newCopyMetrics(scope.NewSubScope("copy")),
		stowMetrics: &stowMetrics{
			BadReference: labeled.NewCounter("bad_key", "Indicates the provided storage reference/key is incorrectly formatted", scope, labeled.EmitUnlabeledMetric),
			BadContainer: labeled.NewCounter("bad_container", "Indicates request for a container that has not been initialized", scope, labeled.EmitUnlabeledMetric),

			HeadFailure: labeled.NewCounter("head_failure", "Indicates failure in HEAD for a given reference", scope, labeled.EmitUnlabeledMetric),
			HeadLatency: labeled.NewStopWatch("head", "Indicates time to fetch metadata using the Head API", time.Millisecond, scope, labeled.EmitUnlabeledMetric),

			ReadFailure:     labeled.NewCounter("read_failure", "Indicates failure in GET for a given reference", scope, labeled.EmitUnlabeledMetric, failureTypeOption),
			ReadOpenLatency: labeled.NewStopWatch("read_open", "Indicates time to first byte when reading", time.Millisecond, scope, labeled.EmitUnlabeledMetric),

			WriteFailure: labeled.NewCounter("write_failure", "Indicates failure in storing/PUT for a given reference", scope, labeled.EmitUnlabeledMetric, failureTypeOption),
			WriteLatency: labeled.NewStopWatch("write", "Time to write an object irrespective of size", time.Millisecond, scope, labeled.EmitUnlabeledMetric),
		},
	}
}

// NewDataStore creates a new Data Store with the supplied config.
func NewDataStore(cfg *Config, metrics *DataStoreMetrics) (s *DataStore, err error) {
	defaultClient := http.DefaultClient
	defer func() {
		http.DefaultClient = defaultClient
	}()

	http.DefaultClient = createHTTPClient(cfg.DefaultHTTPClient)

	var rawStore RawStore
	if fn, found := stores[cfg.Type]; found {
		rawStore, err = fn(cfg, metrics)
		if err != nil {
			return &emptyStore, err
		}

		protoStore := NewDefaultProtobufStore(newCachedRawStore(cfg, rawStore, metrics), metrics)
		return NewCompositeDataStore(NewURLPathConstructor(), protoStore), nil
	}

	return &emptyStore, fmt.Errorf("type is of an invalid value [%v]", cfg.Type)
}

// NewCompositeDataStore composes a new DataStore.
func NewCompositeDataStore(refConstructor ReferenceConstructor, composedProtobufStore ComposedProtobufStore) *DataStore {
	return &DataStore{
		ReferenceConstructor:  refConstructor,
		ComposedProtobufStore: composedProtobufStore,
	}
}
