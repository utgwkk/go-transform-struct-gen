package internal

import "strings"

const StructTagFieldName = "transform_struct"

type StructTag struct {
	DestinationField string
	Skip             bool
}

func ParseStructTag(s string) *StructTag {
	t := &StructTag{}
	if s == "" {
		return t
	}

	if s == "-" {
		return &StructTag{
			Skip: true,
		}
	}

	attrs := strings.Split(s, ",")
	for _, attr := range attrs {
		xs := strings.Split(attr, "=")
		k := xs[0]
		if k == "dst_field" && len(xs) >= 2 {
			v := xs[1]
			t.DestinationField = v
		}
	}

	return t
}
