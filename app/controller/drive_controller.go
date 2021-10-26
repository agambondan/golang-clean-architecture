package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"go-blog-api/app/repository"
//	"go-blog-api/app/security"
//	"go-blog-api/app/service"
//)
//
//type driveController struct {
//	userService service.UserService
//	redis       security.Interface
//	auth        security.TokenInterface
//}
//
//type DriveController interface {
//	CreateFile(c *gin.Context)
//	UpdateFile(c *gin.Context)
//}
//
//func NewDriveController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) DriveController {
//	newDriveService := service.NewUserService(repo.User)
//	return &driveController{newDriveService, redis, auth}
//}
//
//
//func (d *driveController) CreateFile(c *gin.Context) {
//	panic("implement me")
//}
//
//func (d *driveController) UpdateFile(c *gin.Context) {
//	panic("implement me")
//}
