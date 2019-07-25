package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	e := Errorf("Code1", "msg")
	assert.NotNil(t, e)
	assert.Equal(t, "[Code1] msg", e.Error())
}

func TestWrapf(t *testing.T) {
	e := Wrapf("Code1", fmt.Errorf("test error"), "msg")
	assert.NotNil(t, e)
	assert.Equal(t, "[Code1] msg, caused by: test error", e.Error())
}

func TestGetErrorCode(t *testing.T) {
	e := Errorf("Code1", "msg")
	assert.NotNil(t, e)
	code, found := GetErrorCode(e)
	assert.True(t, found)
	assert.Equal(t, "Code1", code)
}

func TestIsCausedBy(t *testing.T) {
	e := Errorf("Code1", "msg")
	assert.NotNil(t, e)

	e = Wrapf("Code2", e, "msg")
	assert.True(t, IsCausedBy(e, "Code1"))
	assert.True(t, IsCausedBy(e, "Code2"))
}

func TestIsCausedByError(t *testing.T) {
	e_root := Errorf("Code1", "msg")
	assert.NotNil(t, e_root)
	e1 := Wrapf("Code2", e_root, "msg")
	assert.True(t, IsCausedByError(e1, e_root))
	e2 := Wrapf("Code3", e1, "msg")
	assert.True(t, IsCausedByError(e2, e_root))
	assert.True(t, IsCausedByError(e2, e1))
}