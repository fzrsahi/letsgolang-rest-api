package dto

type CategoryUpdateDto struct {
	Id   int
	Name string `json:"name" validate:"required"`
}
