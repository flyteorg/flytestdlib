package viper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringToByteArray(t *testing.T) {
	t.Run("Expected types", func(t *testing.T) {
		res, err := stringToByteArray(reflect.TypeOf("hello"), reflect.TypeOf([]byte{}), "hello")
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello"), res)
	})

	t.Run("Unexpected types", func(t *testing.T) {
		input := 5
		res, err := stringToByteArray(reflect.TypeOf(input), reflect.TypeOf([]byte{}), input)
		assert.NoError(t, err)
		assert.NotEqual(t, []byte("hello"), res)
	})
}
