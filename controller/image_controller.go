package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"golang-youtube-api/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type imageController struct {
	userService service.UserService
	postService service.PostService
	redis       security.Interface
	auth        security.TokenInterface
}

type ImageController interface {
	GetImagesByUserId(c *gin.Context)
	GetImagesByPostId(c *gin.Context)
	GetImages(c *gin.Context)
}

func NewImageController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) ImageController {
	newUserService := service.NewUserService(repo.User)
	newPostService := service.NewPostService(repo.Post)
	return &imageController{newUserService, newPostService, redis, auth}
}

func (i *imageController) GetImagesByUserId(ctx *gin.Context) {
	var err error
	var filenames []string
	idParam := ctx.Param("id")
	userId := uuid.MustParse(idParam)
	userFindById, err := i.userService.FindById(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	filenames = strings.Split(userFindById.PhotoProfile, ", ")
	for i := 0; i < len(filenames); i++ {
		buffer := utils.WriteImage(userFindById.UUID, "user", filenames[i])
		ctx.Writer.Header().Set("X-Frame-Options", "DENY")
		ctx.Writer.Header().Set("Vary", "Origin")
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("Referrer-Policy", "same-origin")
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename:%s", filenames[i]))
		ctx.Writer.Header().Set("Content-Type", "image/jpeg")
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := ctx.Writer.Write(buffer.Bytes()); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "unable to write image."})
			log.Println("unable to write image.")
			return
		}
	}
}

func (i *imageController) GetImagesByPostId(ctx *gin.Context) {
	var err error
	var filenames []string
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	postFindById, err := i.postService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	filenames = strings.Split(postFindById.Thumbnail, ", ")
	for i := 0; i < len(filenames); i++ {
		buffer := utils.WriteImage(postFindById.UserUUID, "post", filenames[i])
		ctx.Writer.Header().Set("X-Frame-Options", "DENY")
		ctx.Writer.Header().Set("Vary", "Origin")
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("Referrer-Policy", "same-origin")
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename:%s", filenames[i]))
		ctx.Writer.Header().Set("Content-Type", "image/jpeg")
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := ctx.Writer.Write(buffer.Bytes()); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "unable to write image."})
			log.Println("unable to write image.")
			return
		}
	}
}

func (i *imageController) GetImages(ctx *gin.Context) {
	var id int
	var err error
	var filenames []string
	var folderName string
	folderName = "user"
	idParam := ctx.Param("id")
	uuidParam := ctx.Param("uuid")
	userId := uuid.MustParse(uuidParam)
	if idParam != "" {
		folderName= "post"
		id, err = strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}
	if idParam == "" && folderName == "user" {
		userFindById, err := i.userService.FindById(userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		filenames = strings.Split(userFindById.PhotoProfile, ", ")
	} else if folderName == "post" {
		postFindById, err := i.postService.FindById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		filenames = strings.Split(postFindById.Thumbnail, ", ")
	} else {
		err = errors.New("can't find image")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(filenames); i++ {
		buffer := utils.WriteImage(userId, folderName, filenames[i])
		ctx.Writer.Header().Set("X-Frame-Options", "DENY")
		ctx.Writer.Header().Set("Vary", "Origin")
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("Referrer-Policy", "same-origin")
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename:%s", filenames[i]))
		ctx.Writer.Header().Set("Content-Type", "image/jpeg")
		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := ctx.Writer.Write(buffer.Bytes()); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "unable to write image."})
			log.Println("unable to write image.")
			return
		}
	}
}

//func (i *imageController) GetImagesByUserId(ctx *gin.Context) {
//	var err error
//	var filenames, htmlEncodeBase64, encodeBase64 []string
//	idParam := ctx.Param("id")
//	userId := uuid.MustParse(idParam)
//	userFindById, err := i.userService.FindById(userId)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, err)
//		return
//	}
//	filenames = strings.Split(userFindById.PhotoProfile, ", ")
//	for i := 0; i < len(filenames); i++ {
//		getImageBase64 := utils.GetImageToBase64(userId, "user", filenames[i])
//		img2html := "<html><body><img src=\"" + getImageBase64 + "\" /></body></html>"
//		htmlEncodeBase64 = append(htmlEncodeBase64, img2html)
//		encodeBase64 = append(encodeBase64, getImageBase64)
//		//ctx.Writer.Write([]byte(fmt.Sprintf(img2html)))
//	}
//	ctx.JSON(http.StatusOK, gin.H{"image_base64": encodeBase64})
//}
//
//func (i *imageController) GetImagesByPostId(ctx *gin.Context) {
//	var err error
//	var filenames, htmlEncodeBase64, encodeBase64 []string
//	idParam := ctx.Param("id")
//	id, err := strconv.Atoi(idParam)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, err)
//		return
//	}
//	postFindById, err := i.postService.FindById(uint64(id))
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, err)
//		return
//	}
//	filenames = strings.Split(postFindById.Thumbnail, ", ")
//	for i := 0; i < len(filenames); i++ {
//		getImageBase64 := utils.GetImageToBase64(postFindById.UserUUID, "post", filenames[i])
//		img2html := "<html><body><img src=\"" + getImageBase64 + "\" /></body></html>"
//		htmlEncodeBase64 = append(htmlEncodeBase64, img2html)
//		encodeBase64 = append(encodeBase64, getImageBase64)
//		//ctx.Writer.Write([]byte(fmt.Sprintf(img2html)))
//	}
//	ctx.JSON(http.StatusOK, gin.H{"image_base64": encodeBase64})
//}