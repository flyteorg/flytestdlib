package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/flyteorg/flytestdlib/logger"
	mysql_driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dbNotExists uint16 = 1049

// Produces the DSN (data source name) for mysql connections
// Example: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func getMysqlDsn(ctx context.Context, mysqlConfig MysqlConfig) string {
	// Add reading from a file in the future
	sqlConfig := mysql_driver.Config{
		User:                    mysqlConfig.User,
		Passwd:                  mysqlConfig.Password,
		Net:                     "tcp",
		Addr:                    fmt.Sprintf("%s:%d", mysqlConfig.Host, mysqlConfig.Port),
		DBName:                  mysqlConfig.DbName,
		Collation:               "",
		AllowCleartextPasswords: true,
		MultiStatements:         true,
	}
	return sqlConfig.FormatDSN()
}

func CreateMysqlDbIfNotExists(ctx context.Context, gormConfig *gorm.Config, mysqlConfig MysqlConfig) (*gorm.DB, error) {
	dsn := getMysqlDsn(ctx, mysqlConfig)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		if isMysqlErrorWithCode(err, dbNotExists) {
			logger.Infof(ctx, "Creating database %v", mysqlConfig.DbName)
			withDefaultDB := MysqlConfig{
				Host:     mysqlConfig.Host,
				Port:     mysqlConfig.Port,
				DbName:   "mysql", // should always exist
				User:     mysqlConfig.User,
				Password: mysqlConfig.Password,
			}
			dsn := getMysqlDsn(ctx, withDefaultDB)
			defaultDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
			if err != nil {
				return nil, fmt.Errorf("error creating connection to default database for %s:%d",
					withDefaultDB.Host, withDefaultDB.Port)
			}
			createDBStatement := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", mysqlConfig.DbName)
			result := defaultDB.Exec(createDBStatement)
			if result.Error != nil {
				return nil, result.Error
			}
			// Now that the database should be there, try again.
			return CreateMysqlDbIfNotExists(ctx, gormConfig, mysqlConfig)
		}
		logger.Debugf(ctx, "Error opening MySQL connection %s", err)
		return nil, err
	}
	return db, nil
}

func isMysqlErrorWithCode(err error, code uint16) bool {
	myErr := &mysql_driver.MySQLError{}
	if !errors.As(err, &myErr) {
		// err chain does not contain a MySQLError
		return false
	}

	// MySQLError found in chain and set to code specified
	return myErr.Number == code
}
