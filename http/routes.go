package http

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/config"
	"golang-youtube-api/controller"
	"golang-youtube-api/http/middlewares"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/security/google/oauth"
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

	newRoleController := controller.NewRoleController(repositories, newRedisDB.Auth, newToken)
	newUserController := controller.NewUserController(repositories, newRedisDB.Auth, newToken)
	newCategoryController := controller.NewCategoryController(repositories, newRedisDB.Auth, newToken)
	newPostController := controller.NewPostController(repositories, newRedisDB.Auth, newToken)
	newLoginController := controller.NewLoginController(repositories, newRedisDB.Auth, newToken)
	newImageController := controller.NewImageController(repositories, newRedisDB.Auth, newToken)

	routes := server.Router

	//Home Routing
	routes.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	// OAuth
	routes.GET("/oauth", googleLogin)
	routes.GET("/oauth/google", oauth.HandleGoogleLogin)
	routes.GET("/oauth/google/callback", oauth.CallBackFromGoogle)

	// Auth Login API
	routes.POST("/login", newLoginController.Login)
	routes.POST("/logout", middlewares.AuthMiddleware(), newLoginController.Logout)
	routes.POST("/refresh", middlewares.AuthMiddleware(), newLoginController.Refresh)

	// Images API
	routes.GET("/images/:uuid", newImageController.GetImages)
	routes.GET("/images/:uuid/:id", newImageController.GetImages)
	routes.GET("/images/user/:username", newImageController.GetImagesByUsername)
	routes.GET("/images/post/:title", newImageController.GetImagesByPostTitle)
	routes.GET("/images/category/:name", newImageController.GetImagesByCategoryName)
	//routes.GET("/camera", broadcast)

	// Role API
	routes.POST("/role", middlewares.AuthMiddleware(), newRoleController.SaveRole)
	routes.GET("/roles", middlewares.AuthMiddleware(), newRoleController.GetRoles)
	routes.GET("/role/:id", middlewares.AuthMiddleware(), newRoleController.GetRole)
	routes.PUT("/role/:id", middlewares.AuthMiddleware(), newRoleController.UpdateRole)
	routes.DELETE("/role/:id", middlewares.AuthMiddleware(), newRoleController.DeleteRole)
	routes.GET("/roles/count", middlewares.AuthMiddleware(), newRoleController.CountRoles)

	// Users API
	routes.POST("/user", newUserController.SaveUser)
	routes.GET("/users", newUserController.GetUsers)
	routes.GET("/user/:id", newUserController.GetUser)
	routes.GET("/users/role/:role_id", newUserController.GetUsersByRoleId)
	routes.PUT("/user/:id", newUserController.UpdateUser)
	routes.DELETE("/user/:id", newUserController.DeleteUser)
	routes.GET("/users/count", newUserController.CountUsers)

	// Category API
	routes.POST("/category", newCategoryController.SaveCategory)
	routes.GET("/categories", newCategoryController.GetCategories)
	routes.GET("/category/:id", newCategoryController.GetCategory)
	routes.PUT("/category/:id", newCategoryController.UpdateCategory)
	routes.DELETE("/category/:id", newCategoryController.DeleteCategory)
	routes.GET("/categories/count", newCategoryController.CountCategories)

	// Post API
	routes.POST("/post", newPostController.SavePost)
	routes.GET("/posts", newPostController.GetPosts)
	routes.GET("/post/:id", newPostController.GetPost)
	routes.GET("/posts/uuid/:id", newPostController.GetPostsByUserId)
	routes.GET("/posts/username/:username", newPostController.GetPostsByUsername)
	routes.GET("/posts/category/:name", newPostController.GetPostsByCategoryName)
	routes.PUT("/post/:id", newPostController.UpdatePost)
	routes.DELETE("/post/:id", newPostController.DeletePost)
	routes.GET("/posts/count", newPostController.CountPosts)

	// Slug API
	routes.GET("/slug/user/:username", newUserController.GetUserByUsername)
	routes.GET("/slug/category/:name", newCategoryController.GetCategoryByName)
	routes.GET("/slug/post/:title", newPostController.GetPostByTitle)
}

func googleLogin(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(view.IndexPage))
}