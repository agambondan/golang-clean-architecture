package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog-api/app/lib"
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
	if userCheck != nil {
		if *userCheck.RoleId != 1 {
			ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
			return
		}
		var category model.Category
		var categoryAPI model.CategoryAPI
		contentType := ctx.ContentType()
		if contentType != "application/json" {
			data := ctx.PostForm("data")
			err = json.Unmarshal([]byte(data), &categoryAPI)
			if err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
				return
			}
		} else {
			if err = ctx.ShouldBindJSON(&category); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
				return
			}
		}
		_ = lib.Merge(categoryAPI, &category)
		validate := category.Validate("")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
			return
		}
		_, err = c.categoryService.Create(&category)
		if err != nil {
			ctx.JSON(http.StatusConflict, model.BuildErrorResponse("failed create category", err.Error(), nil))
			return
		}
		fileHeader, _ := ctx.FormFile("images")
		if contentType != "application/json" && fileHeader.Filename != "" {
			uploadFile, err := utils.UploadImageFileToAssets(ctx, "categories", "", utils.DriveCategoriesId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed upload image to google drive", "user failed to created", category))
				return
			}
			category.Image = &uploadFile.Name
			imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
			category.ImageURL = &imageURL
			category.ThumbnailURL = &uploadFile.ThumbnailLink
		}
		c.categoryService.UpdateById(*category.ID, &category)
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", category))
	} else {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	}
}

func (c *categoryController) GetCategories(ctx *gin.Context) {
	limit, offset := utils.GetLimitOffsetParam(ctx)
	categories, err := c.categoryService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("categories not found", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", categories))
}

func (c *categoryController) GetCategory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	categoryFindById, err := c.categoryService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("category not found", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", categoryFindById))
}

func (c *categoryController) GetCategoryByName(ctx *gin.Context) {
	name := ctx.Param("name")
	categoryFindByName, err := c.categoryService.FindByName(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("category not found", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", categoryFindByName))
}

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	category := model.Category{}
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	_, err = c.categoryService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("category not found", err.Error(), nil))
		return
	}
	if err = ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
		return
	}
	_, err = c.categoryService.UpdateById(int64(id), &category)
	if err != nil {
		ctx.JSON(http.StatusConflict, model.BuildErrorResponse("failed to update category", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", category))
}

func (c *categoryController) DeleteCategory(ctx *gin.Context) {
	userCheck, err := utils.AdminAuthMiddleware(c.auth, c.redis, c.userService, c.roleService, ctx, "admin")
	if *userCheck.RoleId != 1 {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), userCheck))
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	err = c.categoryService.DeleteById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed to delete category", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "successfully delete category", nil))
}

func (c *categoryController) CountCategories(ctx *gin.Context) {
	count, err := c.categoryService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed to count category", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", count))
}
