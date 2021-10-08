package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"go-blog-api/app/utils"
	"log"
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
	if contentType != "application/json" {
		firstName := ctx.PostForm("first_name")
		lastName := ctx.PostForm("last_name")
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		username := ctx.PostForm("username")
		phoneNumber := ctx.PostForm("phone_number")
		roleID, _ := strconv.ParseInt(ctx.PostForm("role_id"), 10, 64)
		instagram := ctx.PostForm("instagram")
		facebook := ctx.PostForm("facebook")
		twitter := ctx.PostForm("twitter")
		linkedIn := ctx.PostForm("linkedin")
		user.FirstName = &firstName
		user.LastName = &lastName
		user.Email = &email
		user.Password = &password
		user.Username = &username
		user.PhoneNumber = &phoneNumber
		user.RoleId = &roleID
		user.Instagram = &instagram
		user.Facebook = &facebook
		user.Twitter = &twitter
		user.LinkedIn = &linkedIn
	} else {
		if err = ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "invalid json",
			})
			return
		}
	}
	validate := user.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, validate)
		return
	}
	roleFindById, err := c.roleService.FindById(user.RoleId)
	if err != nil || *roleFindById.Name == "" {
		role := model.Role{}
		roleName := "admin"
		role.Name = &roleName
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
	_, err = c.userService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	uploadFile, err := utils.UploadImageFileToAssets(ctx, "user", user.ID.String(), utils.DriveImagesId)
	if err != nil {
		log.Println(err)
		_ = c.userService.DeleteById(user.ID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": err.Error()})
		return
	}
	user.Image = &uploadFile.Name
	imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
	user.ImageURL = &imageURL
	user.ThumbnailURL = &uploadFile.ThumbnailLink
	_, err = c.userService.UpdateById(user.ID, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": user, "error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user, "filename": uploadFile.Name})
}

func (c *userController) GetUsers(ctx *gin.Context) {
	var usersType model.Users
	limit, offset := utils.GetLimitOffsetParam(ctx)
	users, err := c.userService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_ = lib.Merge(users, &usersType)
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if *userCheck.Role.Name == utils.RoleNameAdmin {
			ctx.JSON(http.StatusOK, users)
			return
		} else {
			ctx.JSON(http.StatusOK, usersType.PublicUsers())
			return
		}
	}
}

func (c *userController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	uuidParam := uuid.MustParse(idParam)
	user, err := c.userService.FindById(&uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if *userCheck.Role.Name == utils.RoleNameAdmin || *userCheck.ID == uuidParam {
			ctx.JSON(http.StatusOK, user)
			return
		} else {
			ctx.JSON(http.StatusOK, user.PublicUser())
			//ctx.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
			return
		}
	}
}

func (c *userController) GetUsersByRoleId(ctx *gin.Context) {
	var usersType model.Users
	idParam := ctx.Param("role_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := c.userService.FindAllByRoleId(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_ = lib.Merge(users, &usersType)
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if *userCheck.Role.Name == utils.RoleNameAdmin {
			ctx.JSON(http.StatusOK, users)
			return
		} else {
			ctx.JSON(http.StatusOK, usersType.PublicUsers())
			return
		}
	}
}

func (c *userController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	findUserByUsername, err := c.userService.FindByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		if *userCheck.Role.Name == utils.RoleNameAdmin {
			ctx.JSON(http.StatusOK, findUserByUsername)
			return
		} else {
			ctx.JSON(http.StatusOK, findUserByUsername.PublicUser())
			return
		}
	}
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var filenames []string
	var user model.User
	idParam := ctx.Param("id")
	uuidParam := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if *checkIdUser.ID != uuidParam || *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't update data, your id not equals"})
		return
	} else {
		contentType := ctx.ContentType()
		if contentType != "application/json" {
			firstName := ctx.PostForm("first_name")
			lastName := ctx.PostForm("last_name")
			email := ctx.PostForm("email")
			password := ctx.PostForm("password")
			username := ctx.PostForm("username")
			phoneNumber := ctx.PostForm("phone_number")
			roleID, _ := strconv.ParseInt(ctx.PostForm("role_id"), 10, 64)
			instagram := ctx.PostForm("instagram")
			facebook := ctx.PostForm("facebook")
			twitter := ctx.PostForm("twitter")
			linkedIn := ctx.PostForm("linkedin")
			user.FirstName = &firstName
			user.LastName = &lastName
			user.Email = &email
			user.Password = &password
			user.Username = &username
			user.PhoneNumber = &phoneNumber
			user.RoleId = &roleID
			user.Instagram = &instagram
			user.Facebook = &facebook
			user.Twitter = &twitter
			user.LinkedIn = &linkedIn
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
		user.ID = checkIdUser.ID
		if contentType != "application/json" {
			filenames, err = utils.CreateUploadPhotoMachine(ctx, user.ID.String(), "/user")
			image := strings.Join(filenames, "")
			user.Image = &image
		}
		userUpdateById, err := c.userService.UpdateById(&uuidParam, &user)
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
	uuidParam := uuid.MustParse(idParam)
	checkIdUser, err := utils.CheckIdUser(c.auth, c.redis, c.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if *checkIdUser.ID != uuidParam && *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "can't update data, your id not equals"})
		return
	} else {
		err := c.userService.DeleteById(&uuidParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete User Successfully"})
	}
}

func (c *userController) CountUsers(ctx *gin.Context) {
	count, err := c.userService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, count)
}
