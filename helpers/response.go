package helpers

type ApiResponse struct {
	StatusCode string      `json:"statusCode"`
	Data       interface{} `json:"data"`
}
