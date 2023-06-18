package internal

import (
	"fmt"
	"strings"
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

	return st, nil
}
