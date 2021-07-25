package main

import (
	"github.com/joho/godotenv"
	"golang-youtube-api/config"
	"golang-youtube-api/http"
	"log"
)

var (
	server      http.Server
	configure   config.Configuration
	pathFileEnv = "/home/agam/IdeaProjects/golang-youtube-api/.env.heroku"
)

func init() {
	if err := godotenv.Load(pathFileEnv); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	server.Run(configure)
}
