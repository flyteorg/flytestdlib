
package storage

import (
	"github.com/graymeta/stow"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
)


func TestIsNotFound(t *testing.T) {
	sysError := &os.PathError{Err: syscall.ENOENT}
	assert.True(t, IsNotFound(sysError))
	flyteError := errors.Wrap(sysError, "Wrapping \"system not found\" error")
	assert.True(t, IsNotFound(flyteError))
	secondLevelError := errors.Wrap(flyteError, "Higher level error")
	assert.True(t, IsNotFound(secondLevelError))

	// more for stow errors
	stowNotFoundError := stow.ErrNotFound
	assert.True(t, IsNotFound(stowNotFoundError))
	flyteError = errors.Wrap(stowNotFoundError, "Wrapping stow.ErrNotFound")
	assert.True(t, IsNotFound(flyteError))
	secondLevelError = errors.Wrap(flyteError, "Higher level error wrapper of the stow.ErrNotFound error")
	assert.True(t, IsNotFound(secondLevelError))
}
