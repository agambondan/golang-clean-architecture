package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang-youtube-api/config"
	"golang-youtube-api/repository"
	"log"
)

type Server struct {
	Router *gin.Engine
}

func (server *Server) Run(config config.Config) {
	newRepositories, err := repository.NewRepositories(config)
	if err != nil {
		log.Fatalln(err)
	}
	server.Router = gin.Default()
	server.routes(newRepositories)
	log.Fatalln(server.Router.Run(":" + config.Port))
}
