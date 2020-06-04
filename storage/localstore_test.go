package storage

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/graymeta/stow/local"

	"github.com/lyft/flytestdlib/contextutils"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/lyft/flytestdlib/promutils/labeled"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalStore(t *testing.T) {
	labeled.SetMetricKeys(contextutils.ProjectKey, contextutils.DomainKey, contextutils.WorkflowIDKey, contextutils.TaskIDKey)
	t.Run("Valid config", func(t *testing.T) {
		testScope := promutils.NewTestScope()
		store, err := newStowRawStore(&Config{
			Stow: &StowConfig{
				Kind: local.Kind,
				Config: map[string]string{
					local.ConfigKeyPath: "./",
				},
			},
			InitContainer: "testdata",
		}, testScope.NewSubScope("x"))

		assert.NoError(t, err)
		assert.NotNil(t, store)

		// Stow local store expects the full path after the container portion (looks like a bug to me)
		rc, err := store.ReadRaw(context.TODO(), DataReference("file://testdata/config.yaml"))
		assert.NoError(t, err)
		if assert.NotNil(t, rc) {
			assert.NoError(t, rc.Close())
		}
	})

	t.Run("Invalid config", func(t *testing.T) {
		testScope := promutils.NewTestScope()
		_, err := newStowRawStore(&Config{}, testScope)
		assert.Error(t, err)
	})

	t.Run("Initialize container", func(t *testing.T) {
		testScope := promutils.NewTestScope()
		tmpDir, err := ioutil.TempDir("", "stdlib_local")
		assert.NoError(t, err)

		stats, err := os.Stat(tmpDir)
		assert.NoError(t, err)
		assert.NotNil(t, stats)

		store, err := newStowRawStore(&Config{
			Stow: &StowConfig{
				Kind: local.Kind,
				Config: map[string]string{
					local.ConfigKeyPath: tmpDir,
				},
			},
			InitContainer: "tmp",
		}, testScope.NewSubScope("y"))

		assert.NoError(t, err)
		assert.NotNil(t, store)

		stats, err = os.Stat(filepath.Join(tmpDir, "tmp"))
		assert.NoError(t, err)
		if assert.NotNil(t, stats) {
			assert.True(t, stats.IsDir())
		}
	})

	t.Run("missing init container", func(t *testing.T) {
		testScope := promutils.NewTestScope()
		tmpDir, err := ioutil.TempDir("", "stdlib_local")
		assert.NoError(t, err)

		stats, err := os.Stat(tmpDir)
		assert.NoError(t, err)
		assert.NotNil(t, stats)

		store, err := newStowRawStore(&Config{
			Stow: &StowConfig{
				Kind: local.Kind,
				Config: map[string]string{
					local.ConfigKeyPath: tmpDir,
				},
			},
		}, testScope.NewSubScope("y"))

		assert.Error(t, err)
		assert.Nil(t, store)
	})

	t.Run("multi-container enabled", func(t *testing.T) {
		testScope := promutils.NewTestScope()
		tmpDir, err := ioutil.TempDir("", "stdlib_local")
		assert.NoError(t, err)

		stats, err := os.Stat(tmpDir)
		assert.NoError(t, err)
		assert.NotNil(t, stats)

		store, err := newStowRawStore(&Config{
			Stow: &StowConfig{
				Kind: local.Kind,
				Config: map[string]string{
					local.ConfigKeyPath: tmpDir,
				},
			},
			InitContainer: "tmp",
			MultiContainerEnabled: true,
		}, testScope.NewSubScope("y"))

		assert.NoError(t, err)
		assert.NotNil(t, store)

		stats, err = os.Stat(filepath.Join(tmpDir, "tmp"))
		assert.NoError(t, err)
		if assert.NotNil(t, stats) {
			assert.True(t, stats.IsDir())
		}
	})
}

