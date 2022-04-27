package remote

import (
	"testing"

	"github.com/flyteorg/flyteadmin/pkg/common"
	"github.com/flyteorg/flytestdlib/promutils"
	"github.com/flyteorg/flytestdlib/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetRemoteDataHandler(t *testing.T) {
	dataStorageClient, err := storage.NewDataStore(storage.GetConfig(), promutils.NewTestScope())
	assert.NotNil(t, err)
	config := RemoteDataHandlerConfig{
		common.GCP,
		1,
		"east-1",
		5,
		"principal",
		dataStorageClient,
	}

	t.Run("gcp config", func(t *testing.T) {
		handler := GetRemoteDataHandler(config)
		assert.NotNil(t, handler)
	})

	t.Run("aws config", func(t *testing.T) {
		config.CloudProvider = common.AWS
		handler := GetRemoteDataHandler(config)
		assert.NotNil(t, handler)
	})

	t.Run("local config", func(t *testing.T) {
		config.CloudProvider = common.Local
		handler := GetRemoteDataHandler(config)
		assert.NotNil(t, handler)
	})

	t.Run("default config", func(t *testing.T) {
		config.CloudProvider = common.None
		handler := GetRemoteDataHandler(config)
		assert.NotNil(t, handler)
	})
}
