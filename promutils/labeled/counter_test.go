package labeled

import (
	"context"
	"strings"
	"testing"

	"github.com/flyteorg/flytestdlib/contextutils"
	"github.com/flyteorg/flytestdlib/promutils"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLabeledCounter(t *testing.T) {
	UnsetMetricKeys()
	assert.NotPanics(t, func() {
		SetMetricKeys(contextutils.ProjectKey, contextutils.DomainKey, contextutils.WorkflowIDKey, contextutils.TaskIDKey)
	})

	t.Run("Labeled", func(t *testing.T) {
		scope := promutils.NewScope("testscope_summary")
		c := NewCounter("c1", "some desc", scope)
		assert.NotNil(t, c)
		ctx := context.TODO()
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		var header = `
			# HELP testscope_summary:c1 some desc
			# TYPE testscope_summary:c1 counter
		`

		var expected = `
			testscope_summary:c1{domain="",project="",task="",wf=""} 2
		`
		err := testutil.CollectAndCompare(c.CounterVec, strings.NewReader(header+expected))
		assert.NoError(t, err)

		ctx = contextutils.WithProjectDomain(ctx, "project", "domain")
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		expected = `
			testscope_summary:c1{domain="",project="",task="",wf=""} 2
			testscope_summary:c1{domain="domain",project="project",task="",wf=""} 2
		`
		err = testutil.CollectAndCompare(c.CounterVec, strings.NewReader(header+expected))
		assert.NoError(t, err)

		ctx = contextutils.WithTaskID(ctx, "task")
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		expected = `
			testscope_summary:c1{domain="",project="",task="",wf=""} 2
			testscope_summary:c1{domain="domain",project="project",task="",wf=""} 2
			testscope_summary:c1{domain="domain",project="project",task="task",wf=""} 2
		`
		err = testutil.CollectAndCompare(c.CounterVec, strings.NewReader(header+expected))
		assert.NoError(t, err)
	})

	t.Run("Unlabeled", func(t *testing.T) {
		scope := promutils.NewScope("testscope_summary")
		c := NewCounter("c2", "some desc", scope, EmitUnlabeledMetricOption{})
		assert.NotNil(t, c)
		ctx := context.TODO()
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		var header = `
			# HELP testscope_summary:c2_unlabeled some desc
			# TYPE testscope_summary:c2_unlabeled counter
		`

		var expected = `
			testscope_summary:c2_unlabeled 2
		`
		err := testutil.CollectAndCompare(c.Counter, strings.NewReader(header+expected))
		assert.NoError(t, err)
	})

	t.Run("AdditionalLabels", func(t *testing.T) {
		scope := promutils.NewScope("testscope_summary")
		opts := AdditionalLabelsOption{Labels: []string{contextutils.ProjectKey.String(), contextutils.ExecIDKey.String()}}
		c := NewCounter("c3", "some desc", scope, opts)
		assert.NotNil(t, c)
		ctx := context.TODO()
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		var header = `
			# HELP testscope_summary:c3 some desc
			# TYPE testscope_summary:c3 counter
		`

		var expected = `
			testscope_summary:c3{domain="",exec_id="",project="",task="",wf=""} 2
		`
		err := testutil.CollectAndCompare(c.CounterVec, strings.NewReader(header+expected))
		assert.NoError(t, err)

		ctx = contextutils.WithExecutionID(ctx, "exec_id")
		c.Inc(ctx)
		c.Add(ctx, 1.0)

		expected = `
			testscope_summary:c3{domain="",exec_id="",project="",task="",wf=""} 2
			testscope_summary:c3{domain="",exec_id="exec_id",project="",task="",wf=""} 2
		`
		err = testutil.CollectAndCompare(c.CounterVec, strings.NewReader(header+expected))
		assert.NoError(t, err)
	})
}
