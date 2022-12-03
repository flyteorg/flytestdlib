package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
)
import "context"
import "testing"

func TestGettingMysqlDsn(t *testing.T) {
	mysql := MysqlConfig{
		Host:     "some.ho.st",
		Port:     3306,
		DbName:   "flyteadmin",
		User:     "user",
		Password: "pass",
	}
	ctx := context.Background()
	xx := getMysqlDsn(ctx, mysql)
	fmt.Printf("connection string!!! %s", xx)
	assert.Equal(t, "user:pass@tcp(some.ho.st:3306)/flyteadmin?allowCleartextPasswords=true&allowNativePasswords=false&checkConnLiveness=false&multiStatements=true&maxAllowedPacket=0", xx)
}
