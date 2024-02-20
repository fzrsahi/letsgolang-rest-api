package product

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"task-one/category"
)

func RegisterRoute(router *httprouter.Router, db *sql.DB) {
	productRepository := NewProductRepository()
	categoryRepository := category.NewCategoryRepository()
	productService := NewProductService(productRepository, db, categoryRepository)
	productController := NewProductController(productService)

	router.GET("/products", productController.FindAll)
	router.GET("/products/:id", productController.FindById)
	router.PATCH("/products/:id", productController.Update)
	router.POST("/products", productController.Create)
	router.DELETE("/products/:id", productController.Delete)
}
