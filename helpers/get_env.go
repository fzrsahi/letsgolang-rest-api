package helpers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

const projectDirName = "task-one" // change to relevant project name

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

type DBConfig struct {
	Connection string
	URI        string
}

type AppConfig struct {
	Port string
}

type Config struct {
	DB        *DBConfig
	AppConfig *AppConfig
}

func GetConfig() *Config {
	loadEnv()

	dbDriver := os.Getenv("DB_DRIVER")
	dbUri := os.Getenv("DB_URI")
	port := os.Getenv("PORT")

	return &Config{
		DB: &DBConfig{
			Connection: dbDriver,
			URI:        dbUri,
		},
		AppConfig: &AppConfig{
			Port: port,
		},
	}
}
