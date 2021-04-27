package files

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	homeDirVal = "/home/user"
	homeDirErr error
)
func TestFindConfigFiles(t *testing.T) {
	t.Run("Find config-* group", func(t *testing.T) {
		files := FindConfigFiles([]string{filepath.Join("testdata", "config*.yaml")})
		assert.Equal(t, 2, len(files))
	})

	t.Run("Find other-group-* group", func(t *testing.T) {
		files := FindConfigFiles([]string{filepath.Join("testdata", "other-group*.yaml")})
		assert.Equal(t, 2, len(files))
	})

	t.Run("Absolute path", func(t *testing.T) {
		files := FindConfigFiles([]string{filepath.Join("testdata", "other-group-1.yaml")})
		assert.Equal(t, 1, len(files))

		files = FindConfigFiles([]string{filepath.Join("testdata", "other-group-3.yaml")})
		assert.Equal(t, 0, len(files))
	})
}

func FakeUserHomeDir() (string, error) {
	return homeDirVal, homeDirErr
}
func TestUserHomeDir(t *testing.T) {
	t.Run("User home dir", func(t *testing.T) {
		osUserHomDir = FakeUserHomeDir
		homeDir := UserHomeDir()
		assert.Equal(t, homeDirVal, homeDir)
	})
	t.Run("User home dir fail", func(t *testing.T) {
		homeDirErr = fmt.Errorf("failed to get users home directory")
		homeDirVal = "."
		osUserHomDir = FakeUserHomeDir
		homeDir := UserHomeDir()
		assert.Equal(t, ".", homeDir)
		// Reset
		homeDirErr = nil
		homeDirVal = "/home/user"
	})
}