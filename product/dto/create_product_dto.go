package dto

type ProductCreateDto struct {
	Name       string `json:"name" validate:"required"`
	CategoryId int    `json:"category_id" validate:"required"`
}
