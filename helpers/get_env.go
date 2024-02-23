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

type MailConfig struct {
	SmtpHost     string
	SmtpPort     int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

type Config struct {
	DB        *DBConfig
	AppConfig *AppConfig
	Redis     *RedisConfig
	Mail      *MailConfig
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

	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	senderName := os.Getenv("CONFIG_SENDER_NAME")
	authEmail := os.Getenv("CONFIG_AUTH_EMAIL")
	authPassword := os.Getenv("CONFIG_AUTH_PASSWORD")

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
		Mail: &MailConfig{
			SmtpHost:     smtpHost,
			SmtpPort:     smtpPort,
			SenderName:   senderName,
			AuthEmail:    authEmail,
			AuthPassword: authPassword,
		},
	}
}
