package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"golang-youtube-api/utils"
	"net/http"
	"strconv"
)

type userController struct {
	userRepo    repository.UserRepository
	userService service.UserService
	redis       security.Interface
	auth        security.TokenInterface
}

type UserController interface {
	SaveUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	GetUsersByRoleId(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController(repo repository.UserRepository, redis security.Interface, auth security.TokenInterface) UserController {
	newUserService := service.NewUserService(repo)
	return &userController{repo, newUserService, redis, auth}
}

func (c *userController) SaveUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	validate := user.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, validate)
		return
	}
	userCreate, err := c.userService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, userCreate)
}

func (c *userController) GetUsers(ctx *gin.Context) {
	var users model.Users
	users, err := c.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	httpStatus, adminCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, ctx, "admin")
	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": adminCheck})
		return
	} else {
		if adminCheck == "admin" {
			ctx.JSON(http.StatusOK, users)
			return
		} else {
			ctx.JSON(http.StatusOK, users.PublicUsers())
			return
		}
	}
}

func (c *userController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	uuid := uuid.MustParse(idParam)
	findById, err := c.userService.FindById(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, findById)

}

func (c *userController) GetUsersByRoleId(ctx *gin.Context) {
	idParam := ctx.Param("role_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	users, err := c.userService.FindAllByRoleId(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	panic("implement me")
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	panic("implement me")
}
