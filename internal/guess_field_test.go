package internal

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadStructType(t *testing.T, packageAndName string) *types.Struct {
	t.Helper()

	s, err := ResolveStruct(packageAndName)
	require.NoError(t, err)

	ty, err := s.LookupType()
	require.NoError(t, err)

	return ty
}

func TestGuessFieldCorrespondings(t *testing.T) {
	testcases := []struct {
		name      string
		dstStruct string
		srcStruct string
		want      map[string]string
	}{
		{
			name:      "without createdAt, updatedAt",
			dstStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/foo.FooModel",
			srcStruct: "*github.com/utgwkk/go-transform-struct-gen/internal/fixtures/bar.BarModel",
			want: map[string]string{
				"Id":   "Id",
				"Name": "Name",
				"Age":  "Age",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dst := loadStructType(t, tc.dstStruct)
			src := loadStructType(t, tc.srcStruct)

			got := GuessFieldCorrespondings(dst, src)
			assert.Equal(t, tc.want, got)
		})
	}
}
