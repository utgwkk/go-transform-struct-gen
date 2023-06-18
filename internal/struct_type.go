package internal

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Struct struct {
	Package     string
	PackageName string
	Name        string
	IsRef       bool
}

func (s *Struct) LookupType() (*types.Struct, error) {
	config := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pp, err := packages.Load(config, s.Package)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %s: %w", s.Package, err)
	}

	pkg := pp[0]
	scope := pkg.Types.Scope()
	baseObj := scope.Lookup(s.Name)
	if baseObj == nil {
		return nil, fmt.Errorf("failed to find %s in package %s", s.Name, s.Package)
	}

	ty := baseObj.Type()
	stTy, ok := ty.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("expected Struct, got %T", ty)
	}

	return stTy, nil
}

func (s *Struct) ReturnTypeString(withoutPackageName bool) string {
	if s.IsRef {
		if withoutPackageName {
			return "*"+s.Name
		}
		return fmt.Sprintf("*%s.%s", s.PackageName, s.Name)
	} else {
		if withoutPackageName {
			return s.Name
		}
		return fmt.Sprintf("%s.%s", s.PackageName, s.Name)
	}
}

func (s *Struct) LiteralTypeString() string {
	if s.IsRef {
		return fmt.Sprintf("&%s.%s", s.PackageName, s.Name)
	} else {
		return fmt.Sprintf("%s.%s", s.PackageName, s.Name)
	}
}
