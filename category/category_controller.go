package category

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"task-one/category/dto"
	"task-one/helpers"
)

type CategoryController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type CategoryControllerImpl struct {
	Service CategoryService
}

func NewCategoryController(categoryService CategoryService) CategoryController {
	return &CategoryControllerImpl{Service: categoryService}
}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryRequest := &dto.CategoryCreateDto{}
	helpers.ReadFromRequestBody(request, categoryRequest)

	data := controller.Service.Create(request.Context(), categoryRequest)
	result := helpers.ApiResponse{
		StatusCode: 201,
		Data:       data,
	}

	helpers.WriteToResponse(writer, result, 201)

}

func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryUpdateRequest := &dto.CategoryUpdateDto{}
	helpers.ReadFromRequestBody(request, categoryUpdateRequest)

	categoryId := params.ByName("id")
	res, err := strconv.Atoi(categoryId)
	helpers.PanicIfError(err)

	categoryUpdateRequest.Id = res

	data := controller.Service.Update(request.Context(), categoryUpdateRequest)
	result := helpers.ApiResponse{
		StatusCode: 201,
		Data:       data,
	}
	helpers.WriteToResponse(writer, result, 201)
}

func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("id")
	res, err := strconv.Atoi(categoryId)
	helpers.PanicIfError(err)

	controller.Service.Delete(request.Context(), res)
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       nil,
	}
	helpers.WriteToResponse(writer, result, 200)

}

func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId := params.ByName("id")
	res, err := strconv.Atoi(categoryId)
	helpers.PanicIfError(err)

	data := controller.Service.FindById(request.Context(), res)
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       data,
	}
	helpers.WriteToResponse(writer, result, 200)

}

func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryResponses := controller.Service.FindAll(request.Context())
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       categoryResponses,
	}
	helpers.WriteToResponse(writer, result, 200)
}
