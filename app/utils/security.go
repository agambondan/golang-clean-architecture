package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-blog-api/app/model"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"strings"
)

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func UpdateUploadPhoto() {

}

func AdminAuthMiddleware(auth security.TokenInterface, redis security.Interface, userService service.UserService, roleService service.RoleService, c *gin.Context, checkRole string) (*model.User, error) {
	var user *model.User
	var role *model.Role
	//check is the user is authenticated first
	metadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		err = errors.New("unauthorized")
		return user, err
	}
	//lookup the metadata in redis:
	userId, err := redis.FetchAuth(metadata.TokenUuid)
	if err != nil {
		err = errors.New("unauthorized")
		return user, err
	}
	if checkRole == "admin" {
		user, err = userService.FindById(&userId)
		if err != nil {
			err = errors.New("user not found")
			return user, err
		}
		role, err = roleService.FindById(user.RoleId)
		if err != nil && err.Error() != "" {
			err = errors.New("role not found")
			return user, err
		}
		user.Role = role
	}
	return user, err
}

func CheckIdUser(auth security.TokenInterface, redis security.Interface, userService service.UserService, c *gin.Context) (*model.User, error) {
	var user *model.User
	var err error
	extractTokenMetadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		return user, err
	}
	userId, err := redis.FetchAuth(extractTokenMetadata.TokenUuid)
	if err != nil {
		return user, err
	}
	user, err = userService.FindById(&userId)
	if err != nil {
		return user, err
	}
	return user, err
}
