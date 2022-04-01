package labeled

import (
	"context"

	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/prometheus/client_golang/prometheus"
)

// Counter represents a counter labeled with values from the context. See labeled.SetMetricsKeys for information about to
// configure that.
type Counter struct {
	*prometheus.CounterVec

	prometheus.Counter
	additionalLabels []contextutils.Key
}

// Inc increments the counter by 1. Use Add to increment it by arbitrary non-negative values. The data point will be
// labeled with values from context. See labeled.SetMetricsKeys for information about how to configure that.
func (c Counter) Inc(ctx context.Context) {
	counter, err := c.CounterVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, c.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	counter.Inc()

	if c.Counter != nil {
		c.Counter.Inc()
	}
}

// Add adds the given value to the counter. It panics if the value is < 0.. The data point will be labeled with values
// from context. See labeled.SetMetricsKeys for information about how to configure that.
func (c Counter) Add(ctx context.Context, v float64) {
	counter, err := c.CounterVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, c.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	counter.Add(v)

	if c.Counter != nil {
		c.Counter.Add(v)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// NewCounter  creates a new labeled counter. Label keys must be set before instantiating a counter. See labeled.SetMetricsKeys for
// information about to configure that.
func NewCounter(name, description string, scope promutils.Scope, opts ...MetricOption) Counter {
	if len(metricKeys) == 0 {
		panic(ErrNeverSet)
	}

	c := Counter{}

	name = promutils.SanitizeMetricName(name)
	for _, opt := range opts {
		if _, emitUnlabeledMetric := opt.(EmitUnlabeledMetricOption); emitUnlabeledMetric {
			c.Counter = scope.MustNewCounter(GetUnlabeledMetricName(name), description)
		} else if additionalLabels, casted := opt.(AdditionalLabelsOption); casted {
			var labels []string
			for _, label := range additionalLabels.Labels {
				if contains(metricStringKeys, label) == false {
					labels = append(labels, label)
				}
			}
			// Here we only append the labels that don't exist in metricStringKeys
			c.CounterVec = scope.MustNewCounterVec(name, description, append(metricStringKeys, labels...)...)
			c.additionalLabels = contextutils.MetricKeysFromStrings(labels)
		}
	}

	if c.CounterVec == nil {
		c.CounterVec = scope.MustNewCounterVec(name, description, metricStringKeys...)
	}

	return c
}
