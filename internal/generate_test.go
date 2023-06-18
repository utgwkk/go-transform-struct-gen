package internal

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name      string
		dstStruct string
		srcStruct string
		opts      *GenerateOption
	}{
		{
			name:      "without createdAt, updatedAt (function)",
			dstStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/foo.FooModel",
			srcStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/bar.BarModel",
			opts: &GenerateOption{
				TransformerType: TransformerTypeFunction,
				TramsformerName: "NewFooFromBar",
				DestinationPath: "./fixtures/bar/to_foo.go",
			},
		},
		{
			name:      "without createdAt, updatedAt (method)",
			dstStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/foo.FooModel",
			srcStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/bar.BarModel",
			opts: &GenerateOption{
				TransformerType: TransformerTypeMethod,
				TramsformerName: "ToFoo",
				DestinationPath: "./fixtures/bar/to_foo.go",
			},
		},
		{
			name:      "without createdAt, updatedAt (method, no ref)",
			dstStruct: "github.com/utgwkk/go-transform-struct-gen/internal/fixtures/foo.FooModel",
			srcStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/bar.BarModel",
			opts: &GenerateOption{
				TransformerType: TransformerTypeMethod,
				TramsformerName: "ToFoo",
				DestinationPath: "./fixtures/bar/to_foo.go",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dst, err := ResolveStruct(tc.dstStruct)
			require.NoError(t, err)
			src, err := ResolveStruct(tc.srcStruct)
			require.NoError(t, err)

			got, err := Generate(dst, src, tc.opts)
			require.NoError(t, err)

			snaps.MatchSnapshot(t, got)
		})
	}
}
