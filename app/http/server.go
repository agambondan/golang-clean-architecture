package http

import (
	"github.com/gin-gonic/gin"
	"go-blog-api/app/config"
	"go-blog-api/app/http/middlewares"
	"go-blog-api/app/repository"
	"log"
)

type Server struct {
	Router *gin.Engine
}

func (server *Server) Run(config config.Configuration) {
	config.Init()
	newRepositories, err := repository.NewRepositories(config)
	if err != nil {
		log.Fatalln(err)
	}
	newRepositories.Migrations()
	newRepositories.Seeder()
	newRepositories.AddForeignKey()
	server.Router = gin.Default()
	server.routes(newRepositories)
	server.Router.Use(middlewares.CORSMiddleware())
	log.Fatalln(server.Router.Run(":" + config.Port))
}
