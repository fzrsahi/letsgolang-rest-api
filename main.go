package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"task-one/configs/route"
	"task-one/helpers"
)

func main() {
	env := helpers.GetConfig()

	PORT := env.AppConfig.Port

	server := http.Server{
		Addr:    PORT,
		Handler: route.Router,
	}

	log.Println("Server Running On : http://" + PORT)

	err := server.ListenAndServe()
	helpers.PanicIfError(err)
}
