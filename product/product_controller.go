package product

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"task-one/helpers"
	"task-one/product/dto"
)

type ProductController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type ProductControllerImpl struct {
	Service ProductService
}

func NewProductController(service ProductService) ProductController {
	return &ProductControllerImpl{Service: service}
}

func (controller *ProductControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productRequest := &dto.ProductCreateDto{}
	helpers.ReadFromRequestBody(request, productRequest)

	data := controller.Service.Create(request.Context(), productRequest)
	result := helpers.ApiResponse{
		StatusCode: 201,
		Data:       data,
	}

	helpers.WriteToResponse(writer, result, 201)
}

func (controller *ProductControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productUpdateRequest := &dto.ProductUpdateDto{}
	helpers.ReadFromRequestBody(request, productUpdateRequest)

	productId := params.ByName("id")
	res, err := strconv.Atoi(productId)
	helpers.PanicIfError(err)

	productUpdateRequest.Id = res

	data := controller.Service.Update(request.Context(), productUpdateRequest)
	result := helpers.ApiResponse{
		StatusCode: 201,
		Data:       data,
	}

	helpers.WriteToResponse(writer, result, 201)

}

func (controller *ProductControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("id")
	res, err := strconv.Atoi(productId)
	helpers.PanicIfError(err)

	controller.Service.Delete(request.Context(), res)
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       nil,
	}
	helpers.WriteToResponse(writer, result, 200)

}

func (controller *ProductControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("id")
	res, err := strconv.Atoi(productId)
	helpers.PanicIfError(err)

	data := controller.Service.FindById(request.Context(), res)
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       data,
	}
	helpers.WriteToResponse(writer, result, 200)

}

func (controller *ProductControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	data := controller.Service.FindAll(request.Context())
	result := helpers.ApiResponse{
		StatusCode: 200,
		Data:       data,
	}
	helpers.WriteToResponse(writer, result, 200)

}
