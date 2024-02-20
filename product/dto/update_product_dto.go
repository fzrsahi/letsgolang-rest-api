package dto

type ProductUpdateDto struct {
	Id         int
	Name       string `json:"name"`
	CategoryId int    `json:"category_id"`
}
