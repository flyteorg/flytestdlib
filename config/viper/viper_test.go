package viper

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringToByteArray(t *testing.T) {
	t.Run("Expected types", func(t *testing.T) {
		input := "hello world"
		base64Encoded := base64.StdEncoding.EncodeToString([]byte(input))
		res, err := stringToByteArray(reflect.TypeOf(base64Encoded), reflect.TypeOf([]byte{}), base64Encoded)
		assert.NoError(t, err)
		assert.Equal(t, []byte(input), res)
	})

	t.Run("Unexpected types", func(t *testing.T) {
		input := 5
		res, err := stringToByteArray(reflect.TypeOf(input), reflect.TypeOf([]byte{}), input)
		assert.NoError(t, err)
		assert.NotEqual(t, []byte("hello"), res)
	})
}
