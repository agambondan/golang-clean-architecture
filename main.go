package main

import (
	"golang-youtube-api/config"
	"golang-youtube-api/http"
)

var (
	server      http.Server
	configure   config.Configuration
	pathFileEnv = "/home/agam/IdeaProjects/golang-youtube-api/.env"
)

func init() {
	//if err := godotenv.Load(pathFileEnv); err != nil {
	//	log.Println("no env gotten")
	//}
}

func main() {
	server.Run(configure)
}
