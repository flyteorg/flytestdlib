package labeled

import (
	"context"

	"github.com/lyft/flytestdlib/contextutils"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/prometheus/client_golang/prometheus"
)

// Represents a counter labeled with values from the context. See labeled.SetMetricsKeys for information about to
// configure that.
type Counter struct {
	*prometheus.CounterVec

	prometheus.Counter
}

// Inc increments the counter by 1. Use Add to increment it by arbitrary non-negative values. The data point will be
// labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (c Counter) Inc(ctx context.Context) {
	counter, err := c.CounterVec.GetMetricWith(contextutils.Values(ctx, metricKeys...))
	if err != nil {
		panic(err.Error())
	}
	counter.Inc()

	if c.Counter != nil {
		c.Counter.Inc()
	}
}

// Add adds the given value to the counter. It panics if the value is < 0.. The data point will be labeled with values
// from context. See labeled.SetMetricsKeys for information about to configure that.
func (c Counter) Add(ctx context.Context, v float64) {
	counter, err := c.CounterVec.GetMetricWith(contextutils.Values(ctx, metricKeys...))
	if err != nil {
		panic(err.Error())
	}
	counter.Add(v)

	if c.Counter != nil {
		c.Counter.Add(v)
	}
}

// Creates a new labeled counter. Label keys must be set before instantiating a counter. See labeled.SetMetricsKeys for
// information about to configure that.
func NewCounter(name, description string, scope promutils.Scope, opts ...MetricOption) Counter {
	if len(metricKeys) == 0 {
		panic(ErrNeverSet)
	}

	c := Counter{
		CounterVec: scope.MustNewCounterVec(name, description, metricStringKeys...),
	}

	for _, opt := range opts {
		if _, emitUnlabeledMetric := opt.(EmitUnlabeledMetricOption); emitUnlabeledMetric {
			c.Counter = scope.MustNewCounter(GetUnlabeledMetricName(name), description)
		}
	}

	return c
}