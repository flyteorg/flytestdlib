package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	errs "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/flyteorg/flytestdlib/ioutils"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/flyteorg/flytestdlib/promutils"
)

type protoMetrics struct {
	FetchLatency                 promutils.StopWatch
	MarshalTime                  promutils.StopWatch
	UnmarshalTime                promutils.StopWatch
	MarshalFailure               prometheus.Counter
	UnmarshalFailure             prometheus.Counter
	WriteFailureUnrelatedToCache prometheus.Counter
	ReadFailureUnrelatedToCache  prometheus.Counter
}

// Implements ProtobufStore to marshal and unmarshal protobufs to/from a RawStore
type DefaultProtobufStore struct {
	RawStore
	metrics *protoMetrics
}

func (s DefaultProtobufStore) ReadProtobuf(ctx context.Context, reference DataReference, msg proto.Message) error {
	rc, err := s.ReadRaw(ctx, reference)
	if err != nil && !IsFailedWriteToCache(err) {
		logger.Errorf(ctx, "Failed to read from the raw store [%s] Error: %v", reference, err)
		s.metrics.ReadFailureUnrelatedToCache.Inc()
		return errs.Wrap(err, fmt.Sprintf("path:%v", reference))
	}

	defer func() {
		err = rc.Close()
		if err != nil {
			logger.Warn(ctx, "Failed to close reference [%v]. Error: %v", reference, err)
		}
	}()

	docContents, err := ioutils.ReadAll(rc, s.metrics.FetchLatency.Start())
	if err != nil {
		return errs.Wrap(err, fmt.Sprintf("readAll: %v", reference))
	}

	t := s.metrics.UnmarshalTime.Start()
	err = proto.Unmarshal(docContents, msg)
	t.Stop()
	if err != nil {
		s.metrics.UnmarshalFailure.Inc()
		return errs.Wrap(err, fmt.Sprintf("unmarshall: %v", reference))
	}

	return nil
}

func (s DefaultProtobufStore) WriteProtobuf(ctx context.Context, reference DataReference, opts Options, msg proto.Message) error {
	t := s.metrics.MarshalTime.Start()
	raw, err := proto.Marshal(msg)
	t.Stop()
	if err != nil {
		s.metrics.MarshalFailure.Inc()
		return err
	}

	err = s.WriteRaw(ctx, reference, int64(len(raw)), opts, bytes.NewReader(raw))
	if err != nil && !IsFailedWriteToCache(err) {
		logger.Errorf(ctx, "Failed to write to the raw store [%s] Error: %v", reference, err)
		s.metrics.WriteFailureUnrelatedToCache.Inc()
		return err
	}
	return nil
}

func NewDefaultProtobufStore(store RawStore, metrics *DataStoreMetrics) DefaultProtobufStore {
	return DefaultProtobufStore{
		RawStore: store,
		metrics:  metrics.protoMetrics,
	}
}
