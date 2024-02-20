package response

type ProductResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
}
