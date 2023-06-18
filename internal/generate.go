package internal

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/tools/imports"
)

type TransformerType string

const (
	TransformerTypeFunction TransformerType = "function"
	TransformerTypeMethod   TransformerType = "method"
)

type GenerateOption struct {
	TransformerType TransformerType
	TramsformerName string
	DestinationPath string
}

func Generate(dst, src *Struct, opts *GenerateOption) (string, error) {
	tDst, err := dst.LookupType()
	if err != nil {
		return "", fmt.Errorf("failed to lookup dst type: %w", err)
	}
	tSrc, err := src.LookupType()
	if err != nil {
		return "", fmt.Errorf("failed to lookup src type: %w", err)
	}

	corr := GuessFieldCorrespondings(tDst, tSrc)

	var sb strings.Builder
	writePackageDecl(&sb, src)
	writeTransformer(&sb, dst, src, opts, corr)

	formatted, err := imports.Process(opts.DestinationPath, []byte(sb.String()), nil)
	if err != nil {
		return "", fmt.Errorf("failed to process import: %w", err)
	}

	return string(formatted), nil
}

func writePackageDecl(sb *strings.Builder, s *Struct) {
	sb.WriteString(fmt.Sprintf("package %s\n\n", s.PackageName))
}

func writeTransformer(sb *strings.Builder, dst, src *Struct, opts *GenerateOption, corr map[string]string) {
	if opts.TransformerType == TransformerTypeFunction {
		sb.WriteString(fmt.Sprintf("func %s(src %s) %s {\n", opts.TramsformerName, src.ReturnTypeString(true), dst.ReturnTypeString(false)))
	} else {
		sb.WriteString(fmt.Sprintf("func (src %s) %s() %s {\n", src.ReturnTypeString(true), opts.TramsformerName, dst.ReturnTypeString(false)))
	}

	sb.WriteString(fmt.Sprintf(" return %s{\n", dst.LiteralTypeString()))

	sortedDstFields := maps.Keys(corr)
	sort.Strings(sortedDstFields)

	for _, dstField := range sortedDstFields {
		srcField := corr[dstField]
		sb.WriteString(fmt.Sprintf("  %s: src.%s,\n", dstField, srcField))
	}

	sb.WriteString(" }\n")
	sb.WriteString("}\n")
}
