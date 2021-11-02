package config

import (
	"bytes"
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

var redisConfig = "mockRedis"
var resourceManagerConfig = ResourceManagerConfig{"mockType", 100, &redisConfig}

type MockAccessor struct {
}

func (MockAccessor) ID() string {
	panic("implement me")
}

func (MockAccessor) InitializeFlags(cmdFlags *flag.FlagSet) {
}

func (MockAccessor) InitializePflags(cmdFlags *pflag.FlagSet) {
}

func (MockAccessor) UpdateConfig(ctx context.Context) error {
	return nil
}

func (MockAccessor) ConfigFilesUsed() []string {
	return []string{"test"}
}

func (MockAccessor) RefreshFromConfig() error {
	return nil
}

func newMockAccessor(options Options) Accessor {
	return MockAccessor{}
}

func executeCommandC(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()

	return buf.String(), err
}

func TestNewConfigCommand(t *testing.T) {
	cmd := NewConfigCommand(newMockAccessor)
	assert.NotNil(t, cmd)

	output, err := executeCommandC(cmd, CommandDiscover)
	assert.NoError(t, err)
	assert.Contains(t, output, "test")

	output, err = executeCommandC(cmd, CommandValidate)
	assert.NoError(t, err)
	assert.Contains(t, output, "test")

	_, err = RegisterSection("resourceManager", &resourceManagerConfig)
	assert.NoError(t, err)
	_, err = executeCommandC(cmd, CommandDocs)
	assert.NoError(t, err)
}

type ResourceManagerConfig struct {
	Type             string  `json:"type" pflag:"noop,Which resource manager to use"`
	ResourceMaxQuota int     `json:"resourceMaxQuota" pflag:",Global limit for concurrent Qubole queries"`
	RedisConfig      *string `json:"" pflag:",Config for Redis resource manager."`
}

func TestGetDefaultValue(t *testing.T) {
	val := getDefaultValue(resourceManagerConfig)
	res := "RedisConfig: mockRedis\nresourceMaxQuota: 100\ntype: mockType\n"
	assert.Equal(t, res, val)
}

func TestGetFieldTypeString(t *testing.T) {
	val := reflect.ValueOf(resourceManagerConfig)
	assert.Equal(t, "config.ResourceManagerConfig", getFieldTypeString(val.Type()))
	assert.Equal(t, "string", getFieldTypeString(val.Field(0).Type()))
	assert.Equal(t, "int", getFieldTypeString(val.Field(1).Type()))
	assert.Equal(t, "string", getFieldTypeString(val.Field(2).Type()))
}

func TestGetFieldDescriptionFromPflag(t *testing.T) {
	val := reflect.ValueOf(resourceManagerConfig)
	assert.Equal(t, "Which resource manager to use", getFieldDescriptionFromPflag(val.Type().Field(0)))
	assert.Equal(t, "Global limit for concurrent Qubole queries", getFieldDescriptionFromPflag(val.Type().Field(1)))
	assert.Equal(t, "Config for Redis resource manager.", getFieldDescriptionFromPflag(val.Type().Field(2)))
}

func TestGetFieldNameFromJSONTag(t *testing.T) {
	val := reflect.ValueOf(resourceManagerConfig)
	assert.Equal(t, "type", getFieldNameFromJSONTag(val.Type().Field(0)))
	assert.Equal(t, "resourceMaxQuota", getFieldNameFromJSONTag(val.Type().Field(1)))
	assert.Equal(t, "RedisConfig", getFieldNameFromJSONTag(val.Type().Field(2)))
}

func TestCanPrint(t *testing.T) {
	assert.True(t, canPrint(resourceManagerConfig))
	assert.True(t, canPrint(&resourceManagerConfig))
	assert.True(t, canPrint([]ResourceManagerConfig{resourceManagerConfig}))
	assert.False(t, canPrint(map[string]ResourceManagerConfig{"config": resourceManagerConfig}))
}
