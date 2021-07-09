package controller

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"net/http"
)

type postController struct {
	postRepo    repository.PostRepository
	postService service.PostService
	userService service.UserService
	redis       security.Interface
	auth        security.TokenInterface
}

type PostController interface {
	SavePost(c *gin.Context)
	GetPosts(c *gin.Context)
	GetPost(c *gin.Context)
	GetPostsByUserId(c *gin.Context)
	GetPostsByCategoryId(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

func NewPostController(postRepo repository.PostRepository, userRepo repository.UserRepository, redis security.Interface, auth security.TokenInterface) PostController {
	newPostService := service.NewPostService(postRepo)
	newUserService := service.NewUserService(userRepo)
	return &postController{postRepo, newPostService, newUserService, redis, auth}
}

func (p *postController) SavePost(ctx *gin.Context) {
	panic("implement me")
}

func (p *postController) GetPosts(ctx *gin.Context) {
	posts := model.Posts{}
	posts, err := p.postService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"posts": posts.PublicPosts()})
}

func (p *postController) GetPost(ctx *gin.Context) {
	panic("implement me")
}

func (p *postController) GetPostsByUserId(ctx *gin.Context) {
	panic("implement me")
}

func (p *postController) GetPostsByCategoryId(ctx *gin.Context) {
	panic("implement me")
}

func (p *postController) UpdatePost(ctx *gin.Context) {
	panic("implement me")
}

func (p *postController) DeletePost(ctx *gin.Context) {
	panic("implement me")
}
