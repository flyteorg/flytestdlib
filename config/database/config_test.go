package database

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/flyteorg/flytestdlib/config"
	"github.com/flyteorg/flytestdlib/config/viper"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/stretchr/testify/assert"
)

func TestParseDatabaseConfig(t *testing.T) {
	assert.NoError(t, logger.SetConfig(&logger.Config{IncludeSourceCode: true}))

	accessor := viper.NewAccessor(config.Options{
		RootSection: Config,
		SearchPaths: []string{filepath.Join("testdata", "config.yaml")},
	})

	assert.NoError(t, accessor.UpdateConfig(context.Background()))

	assert.Equal(t, true, GetConfig().EnableForeignKeyConstraintWhenMigrating)
	assert.Equal(t, 5, GetConfig().MaxOpenConnections)
	assert.Equal(t, 5, GetConfig().MaxIdleConnections)
	assert.Equal(t, config.Duration(config.Duration{Duration: 10000000000}), GetConfig().ConnMaxLifeTime)

	assert.Equal(t, 5432, GetConfig().PostgresConfig.Port)
	assert.Equal(t, "postgres", GetConfig().PostgresConfig.User)
	assert.Equal(t, "localhost", GetConfig().PostgresConfig.Host)
	assert.Equal(t, "postgres", GetConfig().PostgresConfig.DbName)
	assert.Equal(t, "sslmode=disable", GetConfig().PostgresConfig.ExtraOptions)
	assert.Equal(t, "admin.db", GetConfig().SQLiteConfig.File)
}
