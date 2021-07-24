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
	"strings"
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
	GetUsersByRoleId(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) UserController {
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &userController{newUserService, newRoleService, redis, auth}
}

func (c *userController) SaveUser(ctx *gin.Context) {
	var user model.User
	var err error
	firstName := ctx.PostForm("first_name")
	lastName := ctx.PostForm("last_name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	username := ctx.PostForm("username")
	phoneNumber := ctx.PostForm("phone_number")
	roleID := ctx.PostForm("role_id")
	if firstName != "" && lastName != "" && email != "" {
		user.FirstName = firstName
		user.LastName = lastName
		user.Email = email
		user.Password = password
		user.Username = username
		user.PhoneNumber = phoneNumber
		user.RoleId, _ = strconv.ParseUint(roleID, 10, 64)
	} else {
		if err = ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "invalid json",
			})
			return
		}
	}
	roleFindById, err := c.roleService.FindById(user.RoleId)
	if err != nil || roleFindById.Name == "" {
		role := model.Role{}
		role.Prepare()
		role.ID = 1
		role.Name = "admin"
		_, err = c.roleService.Create(&role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}
	roleFindById, err = c.roleService.FindById(user.RoleId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Role not found"})
		return
	}
	validate := user.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, validate)
		return
	}
	user.Prepare()
	filenames, err := utils.CreateUploadPhoto(ctx, user.UUID, "/user")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user.Image = strings.Join(filenames, ", ")
	userCreate, err := c.userService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": userCreate, "filenames": filenames})
}

func (c *userController) GetUsers(ctx *gin.Context) {
	var users model.Users
	users, err := c.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if userCheck.Role.Name == "admin" {
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
	user, err := c.userService.FindById(uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if userCheck.Role.Name == "admin" || userCheck.UUID == uuid {
			ctx.JSON(http.StatusOK, user)
			return
		} else {
			//ctx.JSON(http.StatusOK, user.PublicUser())
			ctx.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
			return
		}
	}
}

func (c *userController) GetUsersByRoleId(ctx *gin.Context) {
	var users model.Users
	idParam := ctx.Param("role_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err = c.userService.FindAllByRoleId(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if userCheck.Role.Name == "admin" {
			ctx.JSON(http.StatusOK, users)
			return
		} else {
			ctx.JSON(http.StatusOK, users.PublicUsers())
			return
		}
	}
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var filenames []string
	var user model.User
	user.Prepare()
	idParam := ctx.Param("id")
	uuid := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if checkIdUser.UUID != uuid && checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't update data, your id not equals"})
		return
	} else {
		contentType := ctx.ContentType()
		if contentType != "application/json" {
			user.FirstName = ctx.PostForm("first_name")
			user.LastName = ctx.PostForm("last_name")
			user.Email = ctx.PostForm("email")
			user.Password = ctx.PostForm("password")
			user.Username = ctx.PostForm("username")
			user.PhoneNumber = ctx.PostForm("phone_number")
			roleID := ctx.PostForm("role_id")
			user.RoleId, _ = strconv.ParseUint(roleID, 10, 64)
		} else {
			if err = ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid json",
				})
				return
			}
		}
		validate := user.Validate("update")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, validate)
			return
		}
		user.UUID = checkIdUser.UUID
		if contentType != "application/json" {
			filenames, err = utils.CreateUploadPhoto(ctx, user.UUID, "/user")
			user.Image = strings.Join(filenames, "")
		}
		userUpdateById, err := c.userService.UpdateById(uuid, &user)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userUpdateById.CreatedAt = checkIdUser.CreatedAt
		ctx.JSON(http.StatusOK, gin.H{"data": userUpdateById, "filenames": filenames})
	}
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	uuid := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if checkIdUser.UUID != uuid && checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't update data, your id not equals"})
		return
	} else {
		err := c.userService.DeleteById(uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete User Successfully"})
	}
}
