package main

import (
	"github.com/joho/godotenv"
	"go-blog-api/app/config"
	"go-blog-api/app/http"
	"log"
)

var (
	server      http.Server
	configure   config.Configuration
	pathFileEnv = "./.env"
)

func init() {
	if err := godotenv.Load(pathFileEnv); err != nil {
		log.Println("no env gotten")
		if err = godotenv.Load("./.env.example"); err != nil {
			log.Println("no env gotten")
		}
	}
}

func main() {
	server.Run(configure)
}
