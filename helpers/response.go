package helpers

type ApiResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}
