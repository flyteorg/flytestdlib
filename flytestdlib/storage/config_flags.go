// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package storage

import (
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "type"), defaultConfig.Type, "Sets the type of storage to configure [s3/minio/local/mem].")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "connection.endpoint"), defaultConfig.Connection.Endpoint.String(), "URL for storage client to connect to.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "connection.auth-type"), defaultConfig.Connection.AuthType, "Auth Type to use [iam, accesskey].")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "connection.access-key"), defaultConfig.Connection.AccessKey, "Access key to use. Only required when authtype is set to accesskey.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "connection.secret-key"), defaultConfig.Connection.SecretKey, "Secret to use when accesskey is set.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "connection.region"), defaultConfig.Connection.Region, "Region to connect to.")
	cmdFlags.Bool(fmt.Sprintf("%v%v", prefix, "connection.disable-ssl"), defaultConfig.Connection.DisableSSL, "Disables SSL connection. Should only be used for development.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "container"), defaultConfig.InitContainer, "Initial container to create -if it doesn't exist-.'")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "cache.max_size_mbs"), defaultConfig.Cache.MaxSizeMegabytes, "Maximum size of the cache where the Blob store data is cached in-memory. If not specified or set to 0,  cache is not used")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "cache.target_gc_percent"), defaultConfig.Cache.TargetGCPercent, "Sets the garbage collection target percentage.")
	cmdFlags.Int64(fmt.Sprintf("%v%v", prefix, "limits.maxDownloadMBs"), defaultConfig.Limits.GetLimitMegabytes, "Maximum allowed download size (in MBs) per call.")
	return cmdFlags
}