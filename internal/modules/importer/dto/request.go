package dto

type PathSourceID struct {
	ID int64 `validate:"required,min=1"`
}
