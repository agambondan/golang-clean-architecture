package main

import (
	"github.com/joho/godotenv"
	"golang-youtube-api/config"
	"golang-youtube-api/http"
	"log"
)

var (
	server    http.Server
	configure config.Configuration
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	server.Run(configure)
}
