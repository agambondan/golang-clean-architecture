package controller

import (
	"github.com/gin-gonic/gin"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"go-blog-api/app/utils"
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
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		var role model.Role
		if err := ctx.ShouldBindJSON(&role); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
		validate := role.Validate("")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
			return
		}
		_, err = c.roleService.Create(&role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("failed to create role", err.Error(), role))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", role))
	}
}

func (c *roleController) GetRoles(ctx *gin.Context) {
	var roles *[]model.Role
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		roles, err = c.roleService.FindAll()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("roles not found", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", roles))
		return
	}
}

func (c *roleController) GetRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
			return
		}
		id64 := int64(id)
		roleFindById, err := c.roleService.FindById(&id64)
		if err != nil {
			ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("role not found", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", roleFindById))
		return
	}
}

func (c *roleController) UpdateRole(ctx *gin.Context) {
	role := model.Role{}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
			return
		}
		id64 := int64(id)
		roleFindById, err := c.roleService.FindById(&id64)
		if err != nil {
			ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("role not found", err.Error(), nil))
			return
		}
		if err = ctx.ShouldBindJSON(&role); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
		roleUpdateById, err := c.roleService.UpdateById(&id64, &role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("can't update roles", err.Error(), roleUpdateById))
			return
		}
		roleUpdateById.CreatedAt = roleFindById.CreatedAt
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", roleUpdateById))
		return
	}
}

func (c *roleController) DeleteRole(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
			return
		}
		id64 := int64(id)
		err = c.roleService.DeleteById(&id64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("can't delete role", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "Successfully delete role", id64))
		return
	}
}

func (c *roleController) CountRoles(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	} else {
		count, err := c.roleService.Count()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed to count role", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, count)
	}
}
