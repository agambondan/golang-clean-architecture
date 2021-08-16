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
	CountRoles(c *gin.Context)
}

func NewRoleController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) RoleController {
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &roleController{newUserService, newRoleService, redis, auth}
}

func (c *roleController) SaveRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		role.Prepare()
		createRole, err := c.roleService.Create(&role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, createRole)
	}
}

func (c *roleController) GetRoles(ctx *gin.Context) {
	var roles []model.Role
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
		roles, err = c.roleService.FindAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		ctx.JSON(http.StatusOK, roles)
		return
	}
}

func (c *roleController) GetRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		roleFindById, err := c.roleService.FindById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, roleFindById)
		return
	}
}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	role := model.Role{}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
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
		if err = ctx.ShouldBindJSON(&role); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		role.Prepare()
		roleUpdateById, err := c.roleService.UpdateById(uint64(id), &role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		roleUpdateById.CreatedAt = roleFindById.CreatedAt
		ctx.JSON(http.StatusOK, roleUpdateById)
		return
	}
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		err = c.roleService.DeleteById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully delete role"})
		return
	}
}

func (c *roleController) CountRoles(ctx *gin.Context) {
	count, err := c.roleService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, count)
}