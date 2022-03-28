package database

import (
	"time"

	"github.com/flyteorg/flytestdlib/config"
	"gorm.io/gorm/logger"
)

const (
	database = "database"
	postgres = "postgres"
)

//go:generate pflags DbConfig --default-var=defaultConfig

var defaultConfig = &DbConfig{
	MaxIdleConnections: 10,
	MaxOpenConnections: 1000,
	ConnMaxLifeTime:    config.Duration{Duration: time.Hour},
	Postgres: &PostgresConfig{
		Port:         5432,
		User:         postgres,
		Host:         postgres,
		DbName:       postgres,
		ExtraOptions: "sslmode=disable",
	},
}
var configSection = config.MustRegisterSection(database, defaultConfig)

// DbConfig is used to for initiating the database connection with the store that holds registered
// entities (e.g. workflows, tasks, launch plans...)
type DbConfig struct {
	DeprecatedHost         string `json:"host" pflag:",deprecated"`
	DeprecatedPort         int    `json:"port" pflag:",deprecated"`
	DeprecatedDbName       string `json:"dbname" pflag:",deprecated"`
	DeprecatedUser         string `json:"username" pflag:",deprecated"`
	DeprecatedPassword     string `json:"password" pflag:",deprecated"`
	DeprecatedPasswordPath string `json:"passwordPath" pflag:",deprecated"`
	DeprecatedExtraOptions string `json:"options" pflag:",deprecated"`
	DeprecatedDebug        bool   `json:"debug" pflag:",deprecated"`

	EnableForeignKeyConstraintWhenMigrating bool            `json:"enableForeignKeyConstraintWhenMigrating" pflag:",Whether to enable gorm foreign keys when migrating the db"`
	MaxIdleConnections                      int             `json:"maxIdleConnections" pflag:",maxIdleConnections sets the maximum number of connections in the idle connection pool."`
	MaxOpenConnections                      int             `json:"maxOpenConnections" pflag:",maxOpenConnections sets the maximum number of open connections to the database."`
	ConnMaxLifeTime                         config.Duration `json:"connMaxLifeTime" pflag:",sets the maximum amount of time a connection may be reused"`
	LogLevel                                logger.LogLevel `json:"log_level" pflag:"-,"`
	Postgres                                *PostgresConfig `json:"postgres,omitempty"`
	SQLite                                  *SQLiteConfig   `json:"sqlite,omitempty"`
}

// SQLiteConfig can be used to configure
type SQLiteConfig struct {
	File string `json:"file" pflag:",The path to the file (existing or new) where the DB should be created / stored. If existing, then this will be re-used, else a new will be created"`
}

// PostgresConfig includes specific config options for opening a connection to a postgres database.
type PostgresConfig struct {
	Host   string `json:"host" pflag:",The host name of the database server"`
	Port   int    `json:"port" pflag:",The port name of the database server"`
	DbName string `json:"dbname" pflag:",The database name"`
	User   string `json:"username" pflag:",The database user who is connecting to the server."`
	// Either Password or PasswordPath must be set.
	Password     string `json:"password" pflag:",The database password."`
	PasswordPath string `json:"passwordPath" pflag:",Points to the file containing the database password."`
	ExtraOptions string `json:"options" pflag:",See http://gorm.io/docs/connecting_to_the_database.html for available options passed, in addition to the above."`
	Debug        bool   `json:"debug" pflag:" Whether or not to start the database connection with debug mode enabled."`
}

func GetConfig() *DbConfig {
	databaseConfig := configSection.GetConfig().(*DbConfig)
	if databaseConfig.Postgres == nil && databaseConfig.SQLite == nil {
		databaseConfig.Postgres = &PostgresConfig{}
	}
	if len(databaseConfig.DeprecatedHost) > 0 {
		databaseConfig.Postgres.Host = databaseConfig.DeprecatedHost
	}
	if databaseConfig.DeprecatedPort != 0 {
		databaseConfig.Postgres.Port = databaseConfig.DeprecatedPort
	}
	if len(databaseConfig.DeprecatedDbName) > 0 {
		databaseConfig.Postgres.DbName = databaseConfig.DeprecatedDbName
	}
	if len(databaseConfig.DeprecatedUser) > 0 {
		databaseConfig.Postgres.User = databaseConfig.DeprecatedUser
	}
	if len(databaseConfig.DeprecatedPassword) > 0 {
		databaseConfig.Postgres.Password = databaseConfig.DeprecatedPassword
	}
	if len(databaseConfig.DeprecatedPasswordPath) > 0 {
		databaseConfig.Postgres.PasswordPath = databaseConfig.DeprecatedPasswordPath
	}
	if len(databaseConfig.DeprecatedExtraOptions) > 0 {
		databaseConfig.Postgres.ExtraOptions = databaseConfig.DeprecatedExtraOptions
	}
	if databaseConfig.DeprecatedDebug != false {
		databaseConfig.Postgres.Debug = databaseConfig.DeprecatedDebug
	}
	return databaseConfig
}
