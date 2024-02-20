package route

import (
	"github.com/julienschmidt/httprouter"
	"task-one/category"
	"task-one/configs/database"
	"task-one/exception"
)

func NewRouter() *httprouter.Router {
	var Router *httprouter.Router = httprouter.New()

	db := database.ConnectToDb()
	category.RegisterRoute(Router, db)

	Router.PanicHandler = exception.ErrorHandler
	return Router
}
