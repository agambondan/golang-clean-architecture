package controller

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"net/http"
)

type roleController struct {
	roleRepo    repository.RoleRepository
	roleService service.RoleService
	redis           security.Interface
	auth            security.TokenInterface
}

type RoleController interface {
	SaveRole(c *gin.Context)
	GetRoles(c *gin.Context)
	GetRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

func NewRoleController(repo repository.RoleRepository, redis security.Interface, auth security.TokenInterface) RoleController {
	newRoleService := service.NewRoleService(repo)
	return &roleController{repo, newRoleService, redis, auth}
}

func (c *roleController) SaveRole(ctx *gin.Context) {
	panic("implement me")
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
	panic("implement me")
}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	panic("implement me")
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	panic("implement me")
}
