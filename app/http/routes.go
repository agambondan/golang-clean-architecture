package http

import (
	"github.com/gin-gonic/gin"
	"go-blog-api/app/config"
	"go-blog-api/app/controller"
	"go-blog-api/app/http/middlewares"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/security/google/oauth"
	"go-blog-api/app/utils/pages/view"
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
	newArticleController := controller.NewArticleController(repositories, newRedisDB.Auth, newToken)
	newLoginController := controller.NewLoginController(repositories, newRedisDB.Auth, newToken)
	newImageController := controller.NewImageController(repositories, newRedisDB.Auth, newToken)

	routes := server.Router

	//public := routes.Group("/api/v1")
	//Home Routing
	routes.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	// OAuth
	routes.GET("/oauth", googleLogin)
	routes.GET("/oauth/google", oauth.HandleGoogleLogin)
	routes.GET("/oauth/google/callback", oauth.CallBackFromGoogle)

	// Auth Login API
	routes.POST("/login", middlewares.CORSMiddleware(), newLoginController.Login)
	routes.POST("/logout", middlewares.CORSMiddleware(), newLoginController.Logout)
	routes.POST("/refresh", middlewares.AuthMiddleware(), newLoginController.Refresh)
	routes.GET("/verify", middlewares.CORSMiddleware(), newLoginController.Verify)
	routes.GET("/verify/role", middlewares.CORSMiddleware(), newLoginController.VerifyRole)

	// Google Drive API

	// Images API
	routes.GET("/images/:uuid", middlewares.CORSMiddleware(), newImageController.GetImages)
	routes.GET("/images/:uuid/:id", middlewares.CORSMiddleware(), newImageController.GetImages)
	routes.GET("/images/user/:username", middlewares.CORSMiddleware(), newImageController.GetImagesByUsername)
	routes.GET("/images/articles/:title", middlewares.CORSMiddleware(), newImageController.GetImagesByArticleTitle)
	routes.GET("/images/categories/:name", middlewares.CORSMiddleware(), newImageController.GetImagesByCategoryName)
	//routes.GET("/camera", broadcast)

	// Role API
	routes.POST("/roles", middlewares.AuthMiddleware(), newRoleController.SaveRole)
	routes.GET("/roles", middlewares.AuthMiddleware(), newRoleController.GetRoles)
	routes.GET("/roles/:id", middlewares.AuthMiddleware(), newRoleController.GetRole)
	routes.PUT("/roles/:id", middlewares.AuthMiddleware(), newRoleController.UpdateRole)
	routes.DELETE("/roles/:id", middlewares.AuthMiddleware(), newRoleController.DeleteRole)
	routes.GET("/roles/count", middlewares.AuthMiddleware(), newRoleController.CountRoles)

	// Users API
	routes.POST("/users", middlewares.CORSMiddleware(), newUserController.SaveUser)
	routes.GET("/users", middlewares.CORSMiddleware(), newUserController.GetUsers)
	routes.GET("/users/:id", middlewares.CORSMiddleware(), newUserController.GetUser)
	routes.GET("/users/role/:role_id", middlewares.CORSMiddleware(), newUserController.GetUsersByRoleId)
	routes.PUT("/users/:id", middlewares.CORSMiddleware(), newUserController.UpdateUser)
	routes.DELETE("/users/:id", middlewares.CORSMiddleware(), newUserController.DeleteUser)
	routes.GET("/users/count", middlewares.CORSMiddleware(), newUserController.CountUsers)
	routes.GET("/users/username/:username", middlewares.CORSMiddleware(), newUserController.GetUserByUsername)

	// Category API
	routes.POST("/categories", middlewares.CORSMiddleware(), newCategoryController.SaveCategory)
	routes.GET("/categories", middlewares.CORSMiddleware(), newCategoryController.GetCategories)
	routes.GET("/categories/:id", middlewares.CORSMiddleware(), newCategoryController.GetCategory)
	routes.PUT("/categories/:id", middlewares.CORSMiddleware(), newCategoryController.UpdateCategory)
	routes.DELETE("/categories/:id", middlewares.CORSMiddleware(), newCategoryController.DeleteCategory)
	routes.GET("/categories/count", middlewares.CORSMiddleware(), newCategoryController.CountCategories)

	// Article API
	routes.POST("/articles", middlewares.CORSMiddleware(), middlewares.CORSMiddleware(), newArticleController.SaveArticle)
	routes.GET("/articles", middlewares.CORSMiddleware(), newArticleController.GetArticles)
	routes.GET("/articles/:id", middlewares.CORSMiddleware(), newArticleController.GetArticle)
	routes.GET("/articles/uuid/:id", middlewares.CORSMiddleware(), newArticleController.GetArticlesByUserId)
	routes.GET("/articles/username/:username", middlewares.CORSMiddleware(), newArticleController.GetArticlesByUsername)
	routes.GET("/articles/categories/:name", middlewares.CORSMiddleware(), newArticleController.GetArticlesByCategoryName)
	routes.GET("/articles/categories/:name/count", middlewares.CORSMiddleware(), newArticleController.GetCountArticlesByCategoryName)
	routes.PUT("/articles/:id", middlewares.CORSMiddleware(), newArticleController.UpdateArticle)
	routes.DELETE("/articles/:id", middlewares.CORSMiddleware(), newArticleController.DeleteArticle)
	routes.GET("/articles/count", middlewares.CORSMiddleware(), newArticleController.CountArticles)

	// Slug API
	routes.GET("/slug/user/:username", middlewares.CORSMiddleware(), newUserController.GetUserByUsername)
	routes.GET("/slug/categories/:name", middlewares.CORSMiddleware(), newCategoryController.GetCategoryByName)
	routes.GET("/slug/articles/:title", middlewares.CORSMiddleware(), newArticleController.GetArticleByTitle)
}

func googleLogin(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(view.IndexPage))
}
