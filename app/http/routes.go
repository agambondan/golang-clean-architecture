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
	routes.POST("/logout", middlewares.CORSMiddleware(), middlewares.AuthMiddleware(), newLoginController.Logout)
	routes.POST("/refresh", middlewares.AuthMiddleware(), newLoginController.Refresh)
	routes.GET("/check", middlewares.CORSMiddleware(), middlewares.AuthMiddleware(), newLoginController.Check)

	// Images API
	routes.GET("/images/:uuid", newImageController.GetImages)
	routes.GET("/images/:uuid/:id", newImageController.GetImages)
	routes.GET("/images/user/:username", newImageController.GetImagesByUsername)
	routes.GET("/images/post/:title", newImageController.GetImagesByArticleTitle)
	routes.GET("/images/category/:name", newImageController.GetImagesByCategoryName)
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
	routes.PUT("/user/:id", middlewares.CORSMiddleware(), newUserController.UpdateUser)
	routes.DELETE("/users/:id", middlewares.CORSMiddleware(), newUserController.DeleteUser)
	routes.GET("/users/count", middlewares.CORSMiddleware(), newUserController.CountUsers)

	// Category API
	routes.POST("/category", newCategoryController.SaveCategory)
	routes.GET("/categories", newCategoryController.GetCategories)
	routes.GET("/category/:id", newCategoryController.GetCategory)
	routes.PUT("/category/:id", newCategoryController.UpdateCategory)
	routes.DELETE("/category/:id", newCategoryController.DeleteCategory)
	routes.GET("/categories/count", newCategoryController.CountCategories)

	// Article API
	routes.POST("/posts", middlewares.CORSMiddleware(), newArticleController.SaveArticle)
	routes.GET("/posts", newArticleController.GetArticles)
	routes.GET("/post/:id", newArticleController.GetArticle)
	routes.GET("/posts/uuid/:id", newArticleController.GetArticlesByUserId)
	routes.GET("/posts/username/:username", newArticleController.GetArticlesByUsername)
	routes.GET("/posts/category/:name", newArticleController.GetArticlesByCategoryName)
	routes.GET("/posts/category/:name/count", newArticleController.GetCountArticlesByCategoryName)
	routes.PUT("/posts/:id", newArticleController.UpdateArticle)
	routes.DELETE("/posts/:id", newArticleController.DeleteArticle)
	routes.GET("/posts/count", newArticleController.CountArticles)

	// Slug API
	routes.GET("/slug/user/:username", newUserController.GetUserByUsername)
	routes.GET("/slug/category/:name", newCategoryController.GetCategoryByName)
	routes.GET("/slug/post/:title", newArticleController.GetArticleByTitle)
}

func googleLogin(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(view.IndexPage))
}
