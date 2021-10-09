package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"go-blog-api/app/utils"
	"net/http"
	"strconv"
)

type userController struct {
	userService service.UserService
	roleService service.RoleService
	redis       security.Interface
	auth        security.TokenInterface
}

type UserController interface {
	SaveUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	GetUserByUsername(c *gin.Context)
	GetUsersByRoleId(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CountUsers(c *gin.Context)
}

func NewUserController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) UserController {
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &userController{newUserService, newRoleService, redis, auth}
}

func (c *userController) SaveUser(ctx *gin.Context) {
	var user model.User
	var err error
	contentType := ctx.ContentType()
	var userAPI model.UserAPI
	if contentType != "application/json" {
		data := ctx.PostForm("data")
		err = json.Unmarshal([]byte(data), &userAPI)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	} else {
		if err = ctx.ShouldBindJSON(&userAPI); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	}
	_ = lib.Merge(userAPI, &user)
	validate := user.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
		return
	}
	_, err = c.roleService.FindById(user.RoleId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("role not found", err.Error(), nil))
		return
	}
	_, err = c.userService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusConflict, model.BuildErrorResponse("failed to create user", err.Error(), user))
		return
	}
	fileHeader, _ := ctx.FormFile("images")
	if contentType != "application/json" && fileHeader.Filename != "" {
		uploadFile, err := utils.UploadImageFileToAssets(ctx, "user", user.ID.String(), utils.DriveImagesId)
		if err != nil {
			_ = c.userService.DeleteById(user.ID)
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed upload image to google drive", "user failed to created", user))
			return
		}
		user.Image = &uploadFile.Name
		imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
		user.ImageURL = &imageURL
		user.ThumbnailURL = &uploadFile.ThumbnailLink
		_, err = c.userService.UpdateById(user.ID, &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed update image data", err.Error(), user))
			return
		}
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", user))
}

func (c *userController) GetUsers(ctx *gin.Context) {
	var usersType model.Users
	limit, offset := utils.GetLimitOffsetParam(ctx)
	users, err := c.userService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("users not found", err.Error(), nil))
		return
	}
	userCheck, _ := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.Role.Name == utils.RoleNameAdmin {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", users))
		return
	} else {
		_ = lib.Merge(users, &usersType)
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", usersType.PublicUsers()))
		return
	}
}

func (c *userController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	uuidParam := uuid.MustParse(idParam)
	user, err := c.userService.FindById(&uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("user not found", err.Error(), user))
		return
	}
	userCheck, _ := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.Role.Name == utils.RoleNameAdmin || *userCheck.ID == uuidParam {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", user))
		return
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", user.PublicUser()))
		return
	}

}

func (c *userController) GetUsersByRoleId(ctx *gin.Context) {
	var usersType model.Users
	idParam := ctx.Param("role_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("role id must number", err.Error(), nil))
		return
	}
	limit, offset := utils.GetLimitOffsetParam(ctx)
	users, err := c.userService.FindAllByRoleId(int64(id), offset, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse(fmt.Sprintf("users with role id %d not found", id), err.Error(), nil))
		return
	}
	userCheck, _ := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.Role.Name == utils.RoleNameAdmin {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", users))
		return
	} else {
		_ = lib.Merge(users, &usersType)
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", usersType.PublicUsers()))
		return
	}
}

func (c *userController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	findUserByUsername, err := c.userService.FindByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse(fmt.Sprintf("user with username %s not found", username), err.Error(), nil))
		return
	}
	userCheck, _ := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.Role.Name == utils.RoleNameAdmin {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", findUserByUsername))
		return
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", findUserByUsername.PublicUser()))
		return
	}
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var user model.User
	var userAPI model.UserAPI
	idParam := ctx.Param("id")
	uuidParam := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), checkIdUser))
		return
	}
	if *checkIdUser.ID != uuidParam && *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("can't update data, your id not equals", err.Error(), checkIdUser))
		return
	}
	contentType := ctx.ContentType()
	if contentType != "application/json" {
		data := ctx.PostForm("data")
		err = json.Unmarshal([]byte(data), &userAPI)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	} else {
		if err = ctx.ShouldBindJSON(&userAPI); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	}
	_ = lib.Merge(userAPI, &user)
	validate := user.Validate("update")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
		return
	}
	_, err = c.userService.UpdateById(&uuidParam, &user)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("can't update user", err.Error(), user))
		return
	}
	fileHeader, _ := ctx.FormFile("images")
	if contentType != "application/json" && fileHeader.Filename != "" {

	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", user))
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	uuidParam := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), checkIdUser))
		return
	}
	if *checkIdUser.ID != uuidParam && *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("can't update data, your id not equals", err.Error(), checkIdUser))
		return
	} else {
		err = c.userService.DeleteById(&uuidParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id not found", err.Error(), uuidParam))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", nil))
	}
}

func (c *userController) CountUsers(ctx *gin.Context) {
	count, err := c.userService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed to count user", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", count))
}
