package main

import (
	"github.com/joho/godotenv"
	"golang-youtube-api/config"
	"golang-youtube-api/http"
	"log"
	"os"
)

var (
	server    http.Server
	configure config.Config
)

func init() {
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	configure.DBDriver = os.Getenv("DB_DRIVER")
	configure.DBHost = os.Getenv("DB_HOST")
	configure.DBPassword = os.Getenv("DB_PASSWORD")
	configure.DBUser = os.Getenv("DB_USER")
	configure.DBName = os.Getenv("DB_NAME")
	configure.DBPort = os.Getenv("DB_PORT")
	configure.Port = os.Getenv("PORT")
	server.Run(configure)
}
