// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsTestType = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementTestType(t reflect.Kind) bool {
	_, exists := dereferencableKindsTestType[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookTestType(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementTestType(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		raw, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			fmt.Printf("Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

func decode_TestType(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookTestType,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_TestType(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_TestType(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_TestType(val, result))
}

func testDecodeSlice_TestType(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_TestType(vStringSlice, result))
}

func TestTestType_GetPFlagSet(t *testing.T) {
	val := TestType{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestTestType_SetFlags(t *testing.T) {
	actual := TestType{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_str", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("str"); err == nil {
				assert.Equal(t, string(DefaultTestType.StringValue), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("str", testValue)
			if vString, err := cmdFlags.GetString("str"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StringValue)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_bl", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vBool, err := cmdFlags.GetBool("bl"); err == nil {
				assert.Equal(t, bool(DefaultTestType.BoolValue), vBool)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("bl", testValue)
			if vBool, err := cmdFlags.GetBool("bl"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vBool), &actual.BoolValue)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_nested.i", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("nested.i"); err == nil {
				assert.Equal(t, int(DefaultTestType.NestedType.IntValue), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("nested.i", testValue)
			if vInt, err := cmdFlags.GetInt("nested.i"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vInt), &actual.NestedType.IntValue)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_ints", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vIntSlice, err := cmdFlags.GetIntSlice("ints"); err == nil {
				assert.Equal(t, []int([]int{12, 1}), vIntSlice)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := join_TestType([]int{12, 1}, ",")

			cmdFlags.Set("ints", testValue)
			if vIntSlice, err := cmdFlags.GetIntSlice("ints"); err == nil {
				testDecodeSlice_TestType(t, join_TestType(vIntSlice, ","), &actual.IntArray)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_strs", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vStringSlice, err := cmdFlags.GetStringSlice("strs"); err == nil {
				assert.Equal(t, []string([]string{"12", "1"}), vStringSlice)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := join_TestType([]string{"12", "1"}, ",")

			cmdFlags.Set("strs", testValue)
			if vStringSlice, err := cmdFlags.GetStringSlice("strs"); err == nil {
				testDecodeSlice_TestType(t, join_TestType(vStringSlice, ","), &actual.StringArray)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_complexArr", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vStringSlice, err := cmdFlags.GetStringSlice("complexArr"); err == nil {
				assert.Equal(t, []string([]string{}), vStringSlice)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1,1"

			cmdFlags.Set("complexArr", testValue)
			if vStringSlice, err := cmdFlags.GetStringSlice("complexArr"); err == nil {
				testDecodeSlice_TestType(t, vStringSlice, &actual.ComplexJSONArray)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_c", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("c"); err == nil {
				assert.Equal(t, string(DefaultTestType.mustMarshalJSON(DefaultTestType.StringToJSON)), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := DefaultTestType.mustMarshalJSON(DefaultTestType.StringToJSON)

			cmdFlags.Set("c", testValue)
			if vString, err := cmdFlags.GetString("c"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StringToJSON)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.type", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.type"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Type), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.type", testValue)
			if vString, err := cmdFlags.GetString("storage.type"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Type)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.endpoint", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.connection.endpoint"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Connection.Endpoint.String()), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := DefaultTestType.StorageConfig.Connection.Endpoint.String()

			cmdFlags.Set("storage.connection.endpoint", testValue)
			if vString, err := cmdFlags.GetString("storage.connection.endpoint"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Connection.Endpoint)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.auth-type", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.connection.auth-type"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Connection.AuthType), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.connection.auth-type", testValue)
			if vString, err := cmdFlags.GetString("storage.connection.auth-type"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Connection.AuthType)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.access-key", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.connection.access-key"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Connection.AccessKey), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.connection.access-key", testValue)
			if vString, err := cmdFlags.GetString("storage.connection.access-key"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Connection.AccessKey)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.secret-key", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.connection.secret-key"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Connection.SecretKey), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.connection.secret-key", testValue)
			if vString, err := cmdFlags.GetString("storage.connection.secret-key"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Connection.SecretKey)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.region", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.connection.region"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.Connection.Region), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.connection.region", testValue)
			if vString, err := cmdFlags.GetString("storage.connection.region"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.Connection.Region)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.connection.disable-ssl", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vBool, err := cmdFlags.GetBool("storage.connection.disable-ssl"); err == nil {
				assert.Equal(t, bool(DefaultTestType.StorageConfig.Connection.DisableSSL), vBool)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.connection.disable-ssl", testValue)
			if vBool, err := cmdFlags.GetBool("storage.connection.disable-ssl"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vBool), &actual.StorageConfig.Connection.DisableSSL)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.container", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vString, err := cmdFlags.GetString("storage.container"); err == nil {
				assert.Equal(t, string(DefaultTestType.StorageConfig.InitContainer), vString)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.container", testValue)
			if vString, err := cmdFlags.GetString("storage.container"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vString), &actual.StorageConfig.InitContainer)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.cache.max_size_mbs", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("storage.cache.max_size_mbs"); err == nil {
				assert.Equal(t, int(DefaultTestType.StorageConfig.Cache.MaxSizeMegabytes), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.cache.max_size_mbs", testValue)
			if vInt, err := cmdFlags.GetInt("storage.cache.max_size_mbs"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vInt), &actual.StorageConfig.Cache.MaxSizeMegabytes)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.cache.target_gc_percent", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("storage.cache.target_gc_percent"); err == nil {
				assert.Equal(t, int(DefaultTestType.StorageConfig.Cache.TargetGCPercent), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.cache.target_gc_percent", testValue)
			if vInt, err := cmdFlags.GetInt("storage.cache.target_gc_percent"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vInt), &actual.StorageConfig.Cache.TargetGCPercent)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_storage.limits.maxDownloadMBs", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt64, err := cmdFlags.GetInt64("storage.limits.maxDownloadMBs"); err == nil {
				assert.Equal(t, int64(DefaultTestType.StorageConfig.Limits.GetLimitMegabytes), vInt64)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("storage.limits.maxDownloadMBs", testValue)
			if vInt64, err := cmdFlags.GetInt64("storage.limits.maxDownloadMBs"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vInt64), &actual.StorageConfig.Limits.GetLimitMegabytes)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_i", func(t *testing.T) {
		t.Run("DefaultValue", func(t *testing.T) {
			// Test that default value is set properly
			if vInt, err := cmdFlags.GetInt("i"); err == nil {
				assert.Equal(t, int(DefaultTestType.elemValueOrNil(DefaultTestType.IntValue).(int)), vInt)
			} else {
				assert.FailNow(t, err.Error())
			}
		})

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("i", testValue)
			if vInt, err := cmdFlags.GetInt("i"); err == nil {
				testDecodeJson_TestType(t, fmt.Sprintf("%v", vInt), &actual.IntValue)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}