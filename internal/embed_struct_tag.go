package internal

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"golang.org/x/tools/go/ast/astutil"
)

func EmbedStructTagToSrc(dst, src *Struct) error {
	embedFile, err := os.CreateTemp(os.TempDir(), "embed-transform-struct-tag")
	if err != nil {
		return err
	}
	defer os.Remove(embedFile.Name())

	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, src.FilePath, nil, 0)
	if err != nil {
		return err
	}

	st, err := resolveStruct(src, parsed)
	if err != nil {
		return err
	}

	rewrited, err := embedStructTag(dst, src, parsed, st)
	if err != nil {
		return err
	}

	if err := format.Node(embedFile, fset, rewrited); err != nil {
		return err
	}
	embedFile.Seek(0, 0)

	srcFile, err := os.Create(src.FilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err := io.Copy(srcFile, embedFile); err != nil {
		return err
	}

	return nil
}

func resolveStruct(src *Struct, f *ast.File) (*ast.StructType, error) {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		spec := genDecl.Specs[0]
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		if typeSpec.Name.Name != src.Name {
			continue
		}
		if structType, ok := typeSpec.Type.(*ast.StructType); ok {
			return structType, nil
		}
	}
	return nil, fmt.Errorf("%s not found", src.Name)
}

func embedStructTag(dst, src *Struct, file *ast.File, s *ast.StructType) (ast.Node, error) {
	dstTy, err := dst.LookupType()
	if err != nil {
		return nil, err
	}

	dstFields := make([]string, 0, dstTy.NumFields())
	for i := 0; i < dstTy.NumFields(); i++ {
		f := dstTy.Field(i)
		if !f.Exported() {
			continue
		}
		dstFields = append(dstFields, f.Name())
	}

	srcTy, err := src.LookupType()
	if err != nil {
		return nil, err
	}

	var rewrited ast.Node = file
	for i := 0; i < srcTy.NumFields(); i++ {
		sf := s.Fields.List[i]
		if sf.Tag != nil && strings.Contains(sf.Tag.Value, StructTagFieldName) {
			continue
		}

		f := srcTy.Field(i)
		if !f.Exported() {
			continue
		}

		distances := calculateLevenshteinDistance(f.Name(), dstFields)
		ignore := false
		if distances[0].distance > 3 {
			ignore = true
		} else {
			if len(distances) >= 2 && distances[0].distance == distances[1].distance {
				log.Printf("destination field of %s is ambiguous (%s and %s)", f.Name(), distances[0].s, distances[1].s)
			}
		}
		candidate := distances[0].s
		rewrited = astutil.Apply(rewrited, nil, func(c *astutil.Cursor) bool {
			node := c.Node()
			genDecl, ok := node.(*ast.GenDecl)
			if !ok {
				return true
			}
			if genDecl.Tok != token.TYPE {
				return true
			}
			spec := genDecl.Specs[0]
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if typeSpec.Name.Name != src.Name {
				return true
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}
			sf := structType.Fields.List[i]
			if sf.Tag != nil && strings.Contains(sf.Tag.Value, StructTagFieldName) {
				return true
			}
			var tagContent string
			if ignore {
				tagContent = fmt.Sprintf("%s:\"-\"", StructTagFieldName)
			} else {
				tagContent = fmt.Sprintf("%s:\"dst_field=%s\"", StructTagFieldName, candidate)
			}
			if sf.Tag == nil {
				sf.Tag = &ast.BasicLit{}
				sf.Tag.Value = fmt.Sprintf("`%s`", tagContent)
			} else {
				// unquote
				sf.Tag.Value = sf.Tag.Value[1:len(sf.Tag.Value)-1]
				sf.Tag.Value = fmt.Sprintf("`%s %s`", tagContent, sf.Tag.Value)
			}
			c.Replace(node)
			return false
		})
	}

	return rewrited, nil
}

type levenshteinDistance struct {
	s        string
	distance int
}

func calculateLevenshteinDistance(srcField string, dstFields []string) []*levenshteinDistance {
	distances := make([]*levenshteinDistance, len(dstFields))
	for i, f := range dstFields {
		distances[i] = &levenshteinDistance{
			s:        f,
			distance: levenshtein.ComputeDistance(strings.ToLower(srcField), strings.ToLower(f)),
		}
	}
	sort.SliceStable(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})
	return distances
}
