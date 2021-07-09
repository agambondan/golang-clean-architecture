package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"net/http"
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
	return value[adjustedPos:len(value)]
}

func FailOnError(c *gin.Context, httpStatus int, err error) {
	if err != nil {
		c.JSON(httpStatus, err)
		return
	}
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
				if lowerRegex[:2] == "pn" || regex[:2] == "jp" {
					dir := filepath.Join("./assets/images/", userId.String(), pathFolder)
					if dir != "" {
						fmt.Println("didalem dir 1")
						err := os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
						if err != nil {
							_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
							_ = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
						}
					}
				}
				filename := filepath.Join("./assets/images/", userId.String(), pathFolder, basename)
				err := c.SaveUploadedFile(file, filename)
				if err != nil {
					return filenames, err
				}
				filenames = append(filenames, file.Filename)
			} else {
				dir := filepath.Join("./assets/images/", userId.String(), pathFolder)
				fmt.Println("didalem dir 2")
				if dir != "" {
					err := os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
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
			fmt.Println("didalem dir 3")
			err := os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
			if err != nil {
				_ = os.Mkdir("./assets/images/"+userId.String(), os.ModePerm)
				_ = os.Mkdir("./assets/images/"+userId.String()+pathFolder, os.ModePerm)
			}
		}
	}
	return filenames, err
}

func AdminAuthMiddleware(auth security.TokenInterface, redis security.Interface, userService service.UserService, c *gin.Context, checkRole string) (int, string, error) {
	//check is the user is authenticated first
	metadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		return http.StatusUnauthorized, "unauthorized", err
	}
	//lookup the metadata in redis:
	userId, err := redis.FetchAuth(metadata.TokenUuid)
	if err != nil {
		return http.StatusUnauthorized, "unauthorized", err
	}
	if checkRole == "admin" {
		user, err := userService.FindById(userId)
		if err != nil {
			return http.StatusBadRequest, "user not found", err
		}
		if user.RoleId != 1 {
			err = errors.New("unauthorized")
			return http.StatusUnauthorized, "your not admin, unauthorized", err
		}
	}
	return http.StatusOK, "admin", err
}
