package internal

import (
	"go/types"
	"reflect"
)

func GuessFieldCorrespondings(dst, src *types.Struct) map[string]string {
	corr := make(map[string]string)
	for j := 0; j < src.NumFields(); j++ {
		sf := src.Field(j)
		if !sf.Exported() {
			continue
		}

		tag := reflect.StructTag(src.Tag(j))
		parsedTag := ParseStructTag(tag.Get(StructTagFieldName))
		if parsedTag.DestinationField != "" {
			corr[parsedTag.DestinationField] = sf.Name()
			continue
		}
	}

	for i := 0; i < dst.NumFields(); i++ {
		df := dst.Field(i)
		if !df.Exported() {
			continue
		}
		if _, ok := corr[df.Name()]; ok {
			continue
		}

		for j := 0; j < src.NumFields(); j++ {
			sf := src.Field(j)
			if !sf.Exported() {
				continue
			}

			if df.Name() == sf.Name() {
				corr[df.Name()] = sf.Name()
			}
		}
	}
	return corr
}
