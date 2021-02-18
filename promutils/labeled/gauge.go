package labeled

import (
	"context"

	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/prometheus/client_golang/prometheus"
)

// Represents a gauge labeled with values from the context. See labeled.SetMetricsKeys for more information
type Gauge struct {
	*prometheus.GaugeVec

	prometheus.Gauge
	additionalLabels []contextutils.Key
}

// Inc increments the gauge by 1. Use Add to increment by arbitrary values. The data point will be
// labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (g Gauge) Inc(ctx context.Context) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, g.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	gauge.Inc()

	if g.Gauge != nil {
		g.Gauge.Inc()
	}
}

// Add adds the given value to the Gauge. (The value can be negative, resulting in a decrease of the Gauge.)
// The data point will be labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (g Gauge) Add(ctx context.Context, v float64) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, g.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	gauge.Add(v)

	if g.Gauge != nil {
		g.Gauge.Add(v)
	}
}

// Set sets the Gauge to an arbitrary value.
// The data point will be labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (g Gauge) Set(ctx context.Context, v float64) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, g.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	gauge.Set(v)

	if g.Gauge != nil {
		g.Gauge.Set(v)
	}
}

// Dec decrements the level by 1. Use Sub to decrement by arbitrary values. The data point will be
// labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (g Gauge) Dec(ctx context.Context) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, g.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	gauge.Dec()

	if g.Gauge != nil {
		g.Gauge.Dec()
	}
}

// Sub adds the given value to the Gauge. The value can be negative, resulting in an increase of the Gauge.
// The data point will be labeled with values from context. See labeled.SetMetricsKeys for information about to configure that.
func (g Gauge) Sub(ctx context.Context, v float64) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, append(metricKeys, g.additionalLabels...)...))
	if err != nil {
		panic(err.Error())
	}
	gauge.Sub(v)

	if g.Gauge != nil {
		g.Gauge.Sub(v)
	}
}

// SetToCurrentTime sets the Gauge to the current Unix time in seconds.
func (g Gauge) SetToCurrentTime(ctx context.Context) {
	gauge, err := g.GaugeVec.GetMetricWith(contextutils.Values(ctx, metricKeys...))
	if err != nil {
		panic(err.Error())
	}
	gauge.SetToCurrentTime()

	if g.Gauge != nil {
		g.Gauge.SetToCurrentTime()
	}
}

// Creates a new labeled gauge. Label keys must be set before instantiating. If the unlabeled option is given,
// this object will also instantiate and emit another gauge with the given name with an _unlabeled suffix.
// See labeled.SetMetricsKeys for information about to configure that.
func NewGauge(name, description string, scope promutils.Scope, opts ...MetricOption) Gauge {
	if len(metricKeys) == 0 {
		panic(ErrNeverSet)
	}

	g := Gauge{}
	for _, opt := range opts {
		if _, emitUnlabeledMetric := opt.(EmitUnlabeledMetricOption); emitUnlabeledMetric {
			g.Gauge = scope.MustNewGauge(GetUnlabeledMetricName(name), description)
		} else if additionalLabels, casted := opt.(AdditionalLabelsOption); casted {
			g.GaugeVec = scope.MustNewGaugeVec(name, description, append(metricStringKeys, additionalLabels.Labels...)...)
			g.additionalLabels = contextutils.MetricKeysFromStrings(additionalLabels.Labels)
		}
	}

	if g.GaugeVec == nil {
		g.GaugeVec = scope.MustNewGaugeVec(name, description, metricStringKeys...)
	}

	return g
}
