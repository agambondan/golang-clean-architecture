package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"go-blog-api/app/utils"
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
	CountCategories(c *gin.Context)
}

func NewCategoryController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) CategoryController {
	newCategoryService := service.NewCategoryService(repo.Category)
	newUserService := service.NewUserService(repo.User)
	newRoleService := service.NewRoleService(repo.Role)
	return &categoryController{newCategoryService, newUserService, newRoleService, redis, auth}
}

func (c *categoryController) SaveCategory(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		var category model.Category
		contentType := ctx.ContentType()
		if contentType != "application/json" {
			name := ctx.PostForm("name")
			category.Name = &name
		} else {
			if err = ctx.ShouldBindJSON(&category); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid json",
				})
				return
			}
		}
		category.Validate("")
		uploadFile, err := utils.UploadImageFileToAssets(ctx, "categories", "", utils.DriveCategoriesId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		category.Image = &uploadFile.Name
		imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
		category.ImageURL = &imageURL
		category.ThumbnailURL = &uploadFile.ThumbnailLink
		createCategory, err := c.categoryService.Create(&category)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, createCategory)
	}
}

func (c *categoryController) GetCategories(ctx *gin.Context) {
	limit, offset := utils.GetLimitOffsetParam(ctx)
	categories, err := c.categoryService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

func (c *categoryController) GetCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryFindById, err := c.categoryService.FindById(int64(id))
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
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		categoryFindById, err := c.categoryService.FindById(int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		if err = ctx.ShouldBindJSON(&category); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		categoryUpdateById, err := c.categoryService.UpdateById(int64(id), &category)
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
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = c.categoryService.DeleteById(int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully delete category"})
		return
	}
}

func (c *categoryController) CountCategories(ctx *gin.Context) {
	count, err := c.categoryService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, count)
}
