package config

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

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

	t.Run("Discover", func(t *testing.T) {
		output, err := executeCommandC(cmd, CommandDiscover)
		assert.NoError(t, err)
		assert.Contains(t, output, "test")
	})

	t.Run("Validate", func(t *testing.T) {
		output, err := executeCommandC(cmd, CommandValidate)
		assert.NoError(t, err)
		assert.Contains(t, output, "test")
	})

	t.Run("Generate", func(t *testing.T) {
		path, err := ioutil.TempDir(os.TempDir(), "generate-test")
		assert.NoError(t, err)

		output, err := executeCommandC(cmd, CommandGenerate, "-o", path)
		assert.NoError(t, err)
		assert.Contains(t, output, "test")
	})
}

type TestType struct {
	String  string    `json:"string,omitempty"`
	Int     int       `json:"i,omitempty"`
	Recurse *TestType `json:"TestType,omitempty"`
}

func TestCreateMap(t *testing.T) {

	s := NewRootSection()
	assert.NoError(t, s.SetConfig(&TestType{
		String: "Root",
		Int:    0,
	}))

	_, err := s.RegisterSection("sub1", &TestType{
		String: "sub1",
		Int:    1,
	})

	assert.NoError(t, err)

	_, err = s.RegisterSection("sub2", &TestType{
		String: "sub2",
		Int:    1,
		Recurse: &TestType{
			String: "subsub",
			Int:    2,
		},
	})

	expected := map[string]interface{}{
		"string": "Root",
		"sub1": map[string]interface{}{
			"string": "sub1",
			"i":      float64(1),
		},
		"sub2": map[string]interface{}{
			"string": "sub2",
			"i":      float64(1),
			"TestType": map[string]interface{}{
				"string": "subsub",
				"i":      float64(2),
			},
		},
	}

	m, err := createMap(s)
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
}

func Test_toMap(t *testing.T) {
	input := TestType{
		String: "something",
		Int:    4,
	}

	expected := map[string]interface{}{
		"string": "something",
		"i":      float64(4),
	}

	m, err := toMap(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
}
