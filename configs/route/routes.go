package route

import (
	"github.com/julienschmidt/httprouter"
	"task-one/category"
	"task-one/configs/database"
)

var Router *httprouter.Router = httprouter.New()

func init() {
	db := database.ConnectToDb()
	category.RegisterRoute(Router, db)
}
