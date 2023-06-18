package bar

import "github.com/utgwkk/go-transform-struct-gen/internal/fixtures/foo"

func (src *BarModel) ToFoo() *foo.FooModel {
	return &foo.FooModel{
		Age:  src.Age,
		Id:   src.Id,
		Name: src.Name,
	}
}
