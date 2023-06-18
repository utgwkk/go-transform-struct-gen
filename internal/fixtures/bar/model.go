package bar

import "time"

type BarModel struct {
	Id        string
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time

	unexported int
}

type BarModel2 struct {
	ID   string `transform_struct:"dst_field=Id"`
	Name string
	Age  int
}

type BarModel3 struct {
	Id   string
	Name string
	Age  int `transform_struct:"-"`
}
