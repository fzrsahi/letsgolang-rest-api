package helpers

import (
	"task-one/category/model"
	"task-one/category/response"
)

func ToCategoryResponse(category model.Category) response.CategoryResponse {
	return response.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ToCategoryResponses(categories []model.Category) []response.CategoryResponse {
	var categoryResponses []response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, ToCategoryResponse(category))
	}
	return categoryResponses
}
