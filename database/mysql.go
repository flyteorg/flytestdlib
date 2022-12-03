package database

import (
	"context"
	"fmt"
	mysql_driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	// TODO: Need to actually create db if not exists, currently manually creating

	dsn := getMysqlDsn(ctx, mysqlConfig)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	return db, err
}
