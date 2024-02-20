package model

import "task-one/product/response"

type Product struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
	CategoryId   int    `json:"category_id"`
}

func ToProductResponse(product Product) response.ProductResponse {
	return response.ProductResponse{
		Id:           product.Id,
		Name:         product.Name,
		CategoryName: product.CategoryName,
	}
}

func ToProductResponses(products []Product) []response.ProductResponse {
	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return productResponses
}
