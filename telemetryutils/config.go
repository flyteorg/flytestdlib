package telemetryutils

import (
	"context"

	"github.com/flyteorg/flytestdlib/config"
	"github.com/flyteorg/flytestdlib/logger"
)

//go:generate pflags Config --default-var=defaultConfig

type Type = string

const configSectionKey = "telemetry"

var (
	ConfigSection = config.MustRegisterSection(configSectionKey, defaultConfig)
	defaultConfig = &Config{
		FileConfig: FileConfig{
			Enabled:  false,
			Filename: "/tmp/trace.txt",
		},
		JaegerConfig: JaegerConfig{
			Enabled:  false,
			Endpoint: "http://localhost:14268/api/traces",
		},
	}
)

type Config struct {
	FileConfig   FileConfig `json:"file", pflag:",TODO"`
	JaegerConfig JaegerConfig `json:"jaeger", pflag:",TODO"`
}

type FileConfig struct {
	Enabled  bool   `json:"enabled" pflag:",TODO"`
	Filename string `json:"filename" pflag:",TODO"`
}

type JaegerConfig struct {
	Enabled  bool   `json:"enabled" pflag:",TODO"`
	Endpoint string `json:"endpoint" pflag:",TODO"`
}

func GetConfig() *Config {
	if c, ok := ConfigSection.GetConfig().(*Config); ok {
		return c
	}

	logger.Warnf(context.TODO(), "Failed to retrieve config section [%v].", configSectionKey)
	return nil
}
