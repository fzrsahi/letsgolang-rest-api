package product

import (
	"fmt"
	"net/http"
)

type ProductController struct {
	ProductService *ProductService
}

func NewProductController(productService *ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func GetAll(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "This Is Product")
}
