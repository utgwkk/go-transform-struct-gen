package internal

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

func ResolveStruct(packageAndName string) (*Struct, error) {
	st := &Struct{
		IsRef: false,
	}

	refIndex := strings.Index(packageAndName, "*")
	if refIndex != -1 {
		st.IsRef = true
		packageAndName = packageAndName[refIndex+1:]
	}

	dotIndex := strings.LastIndex(packageAndName, ".")
	if dotIndex == -1 {
		return nil, fmt.Errorf("dot (.) not found: %s", packageAndName)
	}

	st.Package = packageAndName[:dotIndex]
	st.Name = packageAndName[dotIndex+1:]

	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedModule | packages.NeedTypes,
	}
	pp, err := packages.Load(config, st.Package)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %s: %w", st.Package, err)
	}

	pkg := pp[0]
	st.PackageName = pkg.Name
	scope := pkg.Types.Scope()
	baseObj := scope.Lookup(st.Name)
	if baseObj == nil {
		return nil, fmt.Errorf("failed to find %s in package %s", st.Name, st.Package)
	}
	pos := pkg.Fset.PositionFor(baseObj.Pos(), true)
	st.FilePath = pos.Filename

	return st, nil
}
