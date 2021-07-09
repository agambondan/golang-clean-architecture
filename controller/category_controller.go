package controller

import (
	"github.com/gin-gonic/gin"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"net/http"
)

type categoryController struct {
	categoryRepo    repository.CategoryRepository
	categoryService service.CategoryService
	redis           security.Interface
	auth            security.TokenInterface
}

type CategoryController interface {
	SaveCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	GetCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

func NewCategoryController(repo repository.CategoryRepository, redis security.Interface, auth security.TokenInterface) CategoryController {
	newCategoryService := service.NewCategoryService(repo)
	return &categoryController{repo, newCategoryService, redis, auth}
}

func (c *categoryController) SaveCategory(ctx *gin.Context) {
	panic("implement me")
}

func (c *categoryController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoryService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (c *categoryController) GetCategory(ctx *gin.Context) {
	panic("implement me")
}

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	panic("implement me")
}

func (c *categoryController) DeleteCategory(ctx *gin.Context) {
	panic("implement me")
}
