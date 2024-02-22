package helpers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
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
	Test_URI   string
}

type RedisConfig struct {
	Host     string
	Password string
	Db       int
}

type AppConfig struct {
	Port string
}

type Config struct {
	DB        *DBConfig
	AppConfig *AppConfig
	Redis     *RedisConfig
}

func GetConfig() *Config {
	loadEnv()

	dbDriver := os.Getenv("DB_DRIVER")
	dbUri := os.Getenv("DB_URI")
	dbTestUri := os.Getenv("DB_TEST_URI")
	port := os.Getenv("PORT")

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &Config{
		DB: &DBConfig{
			Connection: dbDriver,
			URI:        dbUri,
			Test_URI:   dbTestUri,
		},
		AppConfig: &AppConfig{
			Port: port,
		},
		Redis: &RedisConfig{
			Host:     redisHost,
			Password: redisPassword,
			Db:       redisDb,
		},
	}
}
