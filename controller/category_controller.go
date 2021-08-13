package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"golang-youtube-api/utils"
	"golang-youtube-api/utils/google"
	"net/http"
	"strconv"
)

type categoryController struct {
	categoryService service.CategoryService
	userService     service.UserService
	roleService     service.RoleService
	redis           security.Interface
	auth            security.TokenInterface
}

type CategoryController interface {
	SaveCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	GetCategory(c *gin.Context)
	GetCategoryByName(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

func NewCategoryController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) CategoryController {
	newCategoryService := service.NewCategoryService(repo.Category)
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &categoryController{newCategoryService, newUserService, newRoleService, redis, auth}
}

func (c *categoryController) SaveCategory(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		var category model.Category
		category.Prepare()
		contentType := ctx.ContentType()
		if contentType != "application/json" {
			category.Name = ctx.PostForm("name")
		} else {
			if err = ctx.ShouldBindJSON(&category); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid json",
				})
				return
			}
		}
		category.Validate("")
		uploadFile, err := google.UploadImageFileToAssets(ctx, "categories", "", utils.DriveCategoriesId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		category.Image = uploadFile.Name
		category.ImageURL = fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
		category.ThumbnailURL = uploadFile.ThumbnailLink
		createCategory, err := c.categoryService.Create(&category)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, createCategory)
	}
}

func (c *categoryController) GetCategories(ctx *gin.Context) {
	var limit, offset int
	var err error
	var categories model.Categories
	queryParamLimit := ctx.Query("_limit")
	queryParamOffset := ctx.Query("_offset")
	if queryParamLimit != "" {
		limit, err = strconv.Atoi(queryParamLimit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	} else {
		limit = 12
	}
	if queryParamOffset != "" {
		offset, err = strconv.Atoi(queryParamOffset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	}
	categories, err = c.categoryService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if categories == nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "sql: no rows in result set"})
		return
	}
	ctx.JSON(http.StatusOK, categories.Categories())
}

func (c *categoryController) GetCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryFindById, err := c.categoryService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categoryFindById)
	return
}

func (c *categoryController) GetCategoryByName(ctx *gin.Context) {
	name := ctx.Param("name")
	categoryFindById, err := c.categoryService.FindByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categoryFindById)
	return
}

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	category := model.Category{}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		categoryFindById, err := c.categoryService.FindById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		if err = ctx.ShouldBindJSON(&category); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		category.Prepare()
		categoryUpdateById, err := c.categoryService.UpdateById(uint64(id), &category)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		categoryUpdateById.CreatedAt = categoryFindById.CreatedAt
		ctx.JSON(http.StatusOK, categoryUpdateById)
		return
	}
}

func (c *categoryController) DeleteCategory(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = c.categoryService.DeleteById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully delete category"})
		return
	}
}
