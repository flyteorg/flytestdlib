package k8s

import (
	"context"
	"testing"

	"github.com/flyteorg/flytestdlib/config"
	"github.com/flyteorg/flytestdlib/config/viper"
	"github.com/stretchr/testify/assert"
)

var providers = []config.AccessorProvider{viper.NewAccessor}

type ExampleConfig struct {
	RemoteClusterConfig ClusterConfig `json:"remote-cluster-config"`
}

var section = config.MustRegisterSection("example", &ExampleConfig{})

func TestClusterConfig(t *testing.T) {
	exampleConfig := section.GetConfig().(*ExampleConfig)
	clusterConfig := exampleConfig.RemoteClusterConfig

	assert.Equal(t, true, clusterConfig.Enabled)
	assert.Equal(t, "127.0.0.1", clusterConfig.Endpoint)

	auth := clusterConfig.Auth
	assert.Equal(t, "/var/run/credentials/token", auth.TokenPath)
	assert.Equal(t, "/var/run/credentials/cacert", auth.CertPath)
	assert.Equal(t, "file_path", auth.Type)
}

func init() {
	var configAccessor = viper.NewAccessor(config.Options{
		StrictMode:  true,
		SearchPaths: []string{"testdata/config.yaml"},
	})

	configAccessor.UpdateConfig(context.TODO())
}
