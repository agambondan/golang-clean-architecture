package http

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/controller"
	"golang-youtube-api/repository"
	"net/http"
)

func Home(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	//ctx.Writer.Write([]byte(view.IndexPage))
}

func (server *Server) routes(repositories *repository.Repositories) {
	newUserController := controller.NewUserController(repositories.User)

	routes := server.Router

	//Home Routing
	routes.GET("/", Home)

	routes.GET("/users", newUserController.GetUsers)
}
