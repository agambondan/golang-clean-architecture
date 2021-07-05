package controller

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/repository"
	"golang-youtube-api/service"
	"net/http"
)

type controller struct {
	userRepo    repository.UserRepository
	userService service.UserService
}

type UserController interface {
	SaveUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController(repo repository.UserRepository) UserController {
	newUserService := service.NewUserService(repo)
	return &controller{repo, newUserService}
}

func (c controller) SaveUser(ctx *gin.Context) {
	panic("implement me")
}

func (c controller) GetUsers(ctx *gin.Context) {
	users, err := c.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c controller) GetUser(ctx *gin.Context) {
	panic("implement me")
}

func (c controller) UpdateUser(ctx *gin.Context) {
	panic("implement me")
}

func (c controller) DeleteUser(ctx *gin.Context) {
	panic("implement me")
}
