package http

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/config"
	"golang-youtube-api/controller"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/utils/pages/view"
	"log"
	"net/http"
)

var (
	configure config.Configuration
)

func (server *Server) routes(repositories *repository.Repositories) {
	configuration := configure.Init()
	newRedisDB, err := security.NewRedisDB(configuration.RedisHost, configuration.RedisPort, configuration.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}
	newToken := security.NewToken()

	newRoleController := controller.NewRoleController(repositories.Role, newRedisDB.Auth, newToken)
	newUserController := controller.NewUserController(repositories.User, newRedisDB.Auth, newToken)
	newCategoryController := controller.NewCategoryController(repositories.Category, newRedisDB.Auth, newToken)
	newPostController := controller.NewPostController(repositories.Post, repositories.User, newRedisDB.Auth, newToken)
	newLoginController := controller.NewLoginController(repositories.User, newRedisDB.Auth, newToken)

	routes := server.Router

	//Home Routing
	routes.GET("/", Home)

	// Auth Login API
	routes.POST("/login", newLoginController.Login)
	routes.POST("/logout", newLoginController.Logout)
	routes.POST("/refresh", newLoginController.Refresh)

	// Role API
	routes.POST("/role", newRoleController.SaveRole)
	routes.GET("/roles", newRoleController.GetRoles)
	routes.GET("/role/:id", newRoleController.GetRole)
	routes.PUT("/role/:id", newRoleController.UpdateRole)
	routes.DELETE("/role/:id", newRoleController.DeleteRole)

	// Users API
	routes.POST("/user", newUserController.SaveUser)
	routes.GET("/users", newUserController.GetUsers)
	routes.GET("/user/:id", newUserController.GetUser)
	routes.GET("/users/role/:role_id", newUserController.GetUsersByRoleId)
	routes.PUT("/user/:id", newUserController.UpdateUser)
	routes.DELETE("/user/:id", newUserController.DeleteUser)

	// Category API
	routes.POST("/category", newCategoryController.SaveCategory)
	routes.GET("/categories", newCategoryController.GetCategories)
	routes.GET("/category/:id", newCategoryController.GetCategory)
	routes.PUT("/category/:id", newCategoryController.UpdateCategory)
	routes.DELETE("/category/:id", newCategoryController.DeleteCategory)

	// Post API
	routes.POST("/post", newPostController.SavePost)
	routes.GET("/posts", newPostController.GetPosts)
	routes.GET("/post/:id", newPostController.GetPost)
	routes.GET("/posts/user/:id", newPostController.GetPostsByUserId)
	routes.GET("/posts/category/:id", newPostController.GetPostsByCategoryId)
	routes.PUT("/post/:id", newPostController.UpdatePost)
	routes.DELETE("/post/:id", newPostController.DeletePost)
}

func Home(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(view.IndexPage))
}
