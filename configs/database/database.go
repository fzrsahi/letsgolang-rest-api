package database

import (
	"database/sql"
	"log"
	"task-one/helpers"
	"time"
)

func ConnectToDb() *sql.DB {
	env := helpers.GetConfig()

	db, err := sql.Open(env.DB.Connection, env.DB.URI)
	helpers.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)

	log.Println("Database Connected..")

	return db

}
