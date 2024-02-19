package dto

type CategoryCreateDto struct {
	Name string `json:"name" validate:"required"`
}
