package api

import (
	"bytes"
	"context"
	"fmt"
	"go/token"
	"go/types"
	"unicode"

	"golang.org/x/tools/go/packages"

	"github.com/flyteorg/flytestdlib/logger"

	"k8s.io/apimachinery/pkg/util/sets"
)

func camelCase(str string) string {
	if len(str) == 0 {
		return str
	}

	firstRune := bytes.Runes([]byte(str))[0]
	if unicode.IsLower(firstRune) {
		return fmt.Sprintf("%v%v", string(unicode.ToUpper(firstRune)), str[1:])
	}

	return str
}

func isJSONUnmarshaler(t types.Type) bool {
	return implementsAnyOfMethods(t, "UnmarshalJSON")
}

func isJSONMarshaler(t types.Type) bool {
	return implementsAnyOfMethods(t, "MarshalJSON")
}

func isStringer(t types.Type) bool {
	return implementsAnyOfMethods(t, "String")
}

func isPFlagValue(t types.Type) bool {
	return implementsAllOfMethods(t, "String", "Set", "Type")
}

func hasStringConstructor(t *types.Named) bool {
	return t.Obj().Parent().Lookup(fmt.Sprintf("%sString", t.Obj().Name())) != nil
}

func implementsAnyOfMethods(t types.Type, methodNames ...string) (found bool) {
	mset := types.NewMethodSet(t)
	for _, name := range methodNames {
		if mset.Lookup(nil, name) != nil {
			return true
		}
	}

	mset = types.NewMethodSet(types.NewPointer(t))
	for _, name := range methodNames {
		if mset.Lookup(nil, name) != nil {
			return true
		}
	}

	return false
}

func implementsAllOfMethods(t types.Type, methodNames ...string) (found bool) {
	fset := token.NewFileSet()
	mset := types.NewMethodSet(t)
	foundMethods := sets.NewString()
	var loadedPackage *packages.Package
	if asNamed, isNamed := t.(*types.Named); isNamed {
		pkg := asNamed.Obj().Pkg()
		config := &packages.Config{
			Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles,
			Logf: logger.InfofNoCtx,
		}

		loadedPkgs, err := packages.Load(config, pkg.Path())
		if err != nil {
			logger.Errorf(context.Background(), err.Error())
		}

		loadedPackage = loadedPkgs[0]
		fset = loadedPackage.Fset
		//// Resolve package path
		//p := filepath.Clean(filepath.Join(os.Getenv("GOPATH"), pkg.Path()))
		////p = gogenutil.StripGopath(p)
		//logger.InfofNoCtx("Loading package from path [%v]", pkg)
		//
		//if pkg == nil {
		//	return false
		//}
		//files, err := ioutil.ReadDir(p)
		//if err != nil {
		//	logger.Errorf(context.Background(), err.Error())
		//}
		////parent := asNamed.Obj().Parent()
		//for _, name := range loadedPkgs[0].GoFiles {
		//	f, _ := os.Stat(name)
		//	fset.AddFile(name, fset.Base(), int(f.Size()))
		//}
	}
	for _, name := range methodNames {
		if foundMethod := mset.Lookup(loadedPackage.Types, name); foundMethod != nil {
			foundMethods.Insert(name)
			pos := foundMethod.Obj().Pos()
			//fileNames := foundMethod.Obj().Pkg().Scope().Names()

			if _, isNamed := t.(*types.Named); isNamed {
				if t.(*types.Named).Obj().Parent().Contains(pos) {
					return true
				}
			}

			p := fset.Position(pos)
			if p.String() == "" {
				return false
			}
		}
	}

	mset = types.NewMethodSet(types.NewPointer(t))
	for _, name := range methodNames {
		if mset.Lookup(nil, name) != nil {
			foundMethods.Insert(name)
		}
	}

	return foundMethods.Len() == len(methodNames)
}
