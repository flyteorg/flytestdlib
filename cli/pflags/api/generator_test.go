package api

import (
	"context"
	"flag"
	"fmt"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Make sure existing config file(s) parse correctly before overriding them with this flag!
var update = flag.Bool("update", false, "Updates testdata")

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		}

		return reflect.ValueOf(v).Interface()
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func TestElemValueOrNil(t *testing.T) {
	var iPtr *int
	assert.Equal(t, 0, elemValueOrNil(iPtr))
	var sPtr *string
	assert.Equal(t, "", elemValueOrNil(sPtr))
	var i int
	assert.Equal(t, 0, elemValueOrNil(i))
	var s string
	assert.Equal(t, "", elemValueOrNil(s))
	var arr []string
	assert.Equal(t, arr, elemValueOrNil(arr))
}

func TestNewGenerator(t *testing.T) {
	g, err := NewGenerator("github.com/flyteorg/flytestdlib/cli/pflags/api", "TestType", "DefaultTestType", false)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	ctx := context.Background()
	p, err := g.Generate(ctx)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	codeOutput, err := ioutil.TempFile("", "output-*.go")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	defer func() { assert.NoError(t, os.Remove(codeOutput.Name())) }()

	testOutput, err := ioutil.TempFile("", "output-*_test.go")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	defer func() { assert.NoError(t, os.Remove(testOutput.Name())) }()

	assert.NoError(t, p.WriteCodeFile(codeOutput.Name()))
	assert.NoError(t, p.WriteTestFile(testOutput.Name()))

	codeBytes, err := ioutil.ReadFile(codeOutput.Name())
	assert.NoError(t, err)

	testBytes, err := ioutil.ReadFile(testOutput.Name())
	assert.NoError(t, err)

	goldenFilePath := filepath.Join("testdata", "testtype.go")
	goldenTestFilePath := filepath.Join("testdata", "testtype_test.go")
	if *update {
		assert.NoError(t, ioutil.WriteFile(goldenFilePath, codeBytes, os.ModePerm))
		assert.NoError(t, ioutil.WriteFile(goldenTestFilePath, testBytes, os.ModePerm))
	}

	goldenOutput, err := ioutil.ReadFile(filepath.Clean(goldenFilePath))
	assert.NoError(t, err)
	assert.Equal(t, string(goldenOutput), string(codeBytes))

	goldenTestOutput, err := ioutil.ReadFile(filepath.Clean(goldenTestFilePath))
	assert.NoError(t, err)
	assert.Equal(t, string(goldenTestOutput), string(testBytes))
	t.Run("empty package", func(t *testing.T) {
		gen, err := NewGenerator("", "TestType", "DefaultTestType", false)
		assert.Nil(t, err)
		assert.NotNil(t, gen.GetTargetPackage())
	})
}

func TestBuildFieldForMap(t *testing.T) {
	t.Run("supported : StringToString", func(t *testing.T) {
		ctx := context.Background()
		key := types.Typ[types.String]
		elem := types.Typ[types.String]
		typesMap := types.NewMap(key, elem)
		name := "m"
		goName := "StringMap"
		usage := "I'm a map of strings"
		defaultValue := "DefaultValue"
		fieldInfo, err := buildFieldForMap(ctx, typesMap, name, goName, usage, defaultValue, false)
		assert.Nil(t, err)
		assert.NotNil(t, fieldInfo)
		assert.Equal(t, "StringToString", fieldInfo.FlagMethodName)
		assert.Equal(t, defaultValue, fieldInfo.DefaultValue)
	})
	t.Run("unsupported : not a string type map", func(t *testing.T) {
		ctx := context.Background()
		key := types.Typ[types.Bool]
		elem := types.Typ[types.Bool]
		typesMap := types.NewMap(key, elem)
		name := "m"
		goName := "BoolMap"
		usage := "I'm a map of bools"
		defaultValue := ""
		fieldInfo, err := buildFieldForMap(ctx, typesMap, name, goName, usage, defaultValue, false)
		assert.Nil(t, err)
		assert.NotNil(t, fieldInfo)
		assert.Equal(t, "StringToString", fieldInfo.FlagMethodName)
		assert.Equal(t, "nil", fieldInfo.DefaultValue)
	})
	t.Run("unsupported : elem not a basic type", func(t *testing.T) {
		ctx := context.Background()
		key := types.Typ[types.String]
		elem := &types.Interface{}
		typesMap := types.NewMap(key, elem)
		name := "m"
		goName := "InterfaceMap"
		usage := "I'm a map of interface values"
		defaultValue := ""
		fieldInfo, err := buildFieldForMap(ctx, typesMap, name, goName, usage, defaultValue, false)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("map of type [interface{/* incomplete */}] is not supported."+
			" Only basic slices or slices of json-unmarshalable types are supported"), err)
		assert.NotNil(t, fieldInfo)
		assert.Equal(t, "", fieldInfo.FlagMethodName)
		assert.Equal(t, "", fieldInfo.DefaultValue)
	})
	t.Run("supported : StringToFloat64", func(t *testing.T) {
		ctx := context.Background()
		key := types.Typ[types.String]
		elem := types.Typ[types.Float64]
		typesMap := types.NewMap(key, elem)
		name := "m"
		goName := "Float64Map"
		usage := "I'm a map of float64"
		defaultValue := "DefaultValue"
		fieldInfo, err := buildFieldForMap(ctx, typesMap, name, goName, usage, defaultValue, false)
		assert.Nil(t, err)
		assert.NotNil(t, fieldInfo)
		assert.Equal(t, "StringToFloat64", fieldInfo.FlagMethodName)
		assert.Equal(t, defaultValue, fieldInfo.DefaultValue)
	})
}

func TestDiscoverFieldsRecursive(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		ctx := context.Background()
		defaultValueAccessor := "defaultAccessor"
		fieldPath := "field.Path"
		pkg := types.NewPackage("p", "p")
		n1 := types.NewTypeName(token.NoPos, pkg, "T1", nil)
		namedTypes := types.NewNamed(n1, new(types.Struct), nil)
		//namedTypes := types.NewNamed(n1, nil, nil)
		fields, err := discoverFieldsRecursive(ctx, namedTypes, defaultValueAccessor, fieldPath, false)
		assert.Nil(t, err)
		assert.Equal(t, len(fields), 0)
	})
}
