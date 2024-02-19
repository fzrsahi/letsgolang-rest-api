package category

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
)

func RegisterRoute(router *httprouter.Router, db *sql.DB) {

	categoryRepository := NewCategoryRepository()
	categoryService := NewCategoryService(categoryRepository, db)
	categoryController := NewCategoryController(categoryService)

	router.POST("/categories", categoryController.Create)
	router.GET("/categories", categoryController.FindAll)
	router.GET("/categories/:id", categoryController.FindById)
	router.DELETE("/categories/:id", categoryController.Delete)
	router.PATCH("/categories/:id", categoryController.Update)
}
