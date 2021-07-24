package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"os"
	"path/filepath"
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

func CreateUploadPhoto(c *gin.Context, userId uuid.UUID, pathFolder string) ([]string, error) {
	// create folder and upload foto
	var err error
	var filenames []string
	header := c.Request.Header
	if header.Get("Content-Type")[:19] == "multipart/form-data" {
		formUser, err := c.MultipartForm()
		if err != nil {
			return filenames, err
		}
		files := formUser.File["images"]
		for _, file := range files {
			if file.Size != 0 {
				basename := filepath.Base(file.Filename)
				regex := After(basename, ".")
				lowerRegex := strings.ToLower(regex)
				if lowerRegex[:2] == "pn" || lowerRegex[:2] == "jp" {
					dir := filepath.Join("./assets/images/", userId.String(), pathFolder)
					if dir != "" {
						err = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
						if err != nil {
							_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
							_ = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
						}
					}
				}
				filename := filepath.Join("./assets/images/", userId.String(), pathFolder, basename)
				err = c.SaveUploadedFile(file, filename)
				if err != nil {
					return filenames, err
				}
				filenames = append(filenames, file.Filename)
			} else {
				dir := filepath.Join("./assets/images/", userId.String(), pathFolder)
				if dir != "" {
					err = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
					if err != nil {
						_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
						_ = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
					}
				}
			}
		}
	} else {
		dir := filepath.Join("./assets/images/", userId.String(), pathFolder)
		if dir != "" {
			err := os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
			if err != nil {
				_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
				_ = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
			}
		}
	}
	return filenames, err
}

func UpdateUploadPhoto() {

}

func AdminAuthMiddleware(auth security.TokenInterface, redis security.Interface, userService service.UserService, roleService service.RoleService, c *gin.Context, checkRole string) (model.User, error) {
	var user model.User
	var role model.Role
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
		user, err = userService.FindById(userId)
		if err != nil {
			err = errors.New("user not found")
			return user, err
		}
		role, err = roleService.FindById(user.RoleId)
		if err != nil && err.Error() != "" {
			err = errors.New("role not found")
			return user, err
		}
		user.Role = &role
		//if user.RoleId != 1 {
		//	err = errors.New("your not admin, unauthorized")
		//	return user, err
		//}
	}
	return user, err
}

func CheckIdUser(auth security.TokenInterface, redis security.Interface, userService service.UserService, c *gin.Context) (model.User, error) {
	var user model.User
	var err error
	extractTokenMetadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		return user, err
	}
	userId, err := redis.FetchAuth(extractTokenMetadata.TokenUuid)
	if err != nil {
		return user, err
	}
	user, err = userService.FindById(userId)
	if err != nil {
		return user, err
	}
	return user, err
}
