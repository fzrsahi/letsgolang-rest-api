package exception

import (
	"net/http"
	"task-one/helpers"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)

}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")

		apiResponse := helpers.ApiResponse{
			StatusCode: http.StatusNotFound,
			Data:       exception.Error,
		}

		helpers.WriteToResponse(writer, apiResponse, http.StatusNotFound)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")

	apiResponse := helpers.ApiResponse{
		StatusCode: http.StatusInternalServerError,
		Data:       nil,
	}

	panic(err)
	helpers.WriteToResponse(writer, apiResponse, http.StatusInternalServerError)
}
