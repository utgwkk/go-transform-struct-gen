package internal

import (
	"go/types"
)

func GuessFieldCorrespondings(dst, src *types.Struct) map[string]string {
	corr := make(map[string]string)
	for i := 0; i < dst.NumFields(); i++ {
		df := dst.Field(i)
		if !df.Exported() {
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
