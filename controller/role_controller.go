package controller

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"golang-youtube-api/utils"
	"net/http"
	"strconv"
)

type roleController struct {
	userService service.UserService
	roleService service.RoleService
	redis       security.Interface
	auth        security.TokenInterface
}

type RoleController interface {
	SaveRole(c *gin.Context)
	GetRoles(c *gin.Context)
	GetRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

func NewRoleController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) RoleController {
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &roleController{newUserService, newRoleService, redis, auth}
}

func (c *roleController) SaveRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.Role.Name != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
		var role model.Role
		if err := ctx.ShouldBindJSON(&role); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid json"})
			return
		}
		validate := role.Validate("")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, validate)
			return
		}
		createRole, err := c.roleService.Create(&role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, createRole)
	}
}

func (c *roleController) GetRoles(ctx *gin.Context) {
	roles, err := c.roleService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (c *roleController) GetRole(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	roleFindById, err := c.roleService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, roleFindById)
}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	panic("implement me")
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	panic("implement me")
}
