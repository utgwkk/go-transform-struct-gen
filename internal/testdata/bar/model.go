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
