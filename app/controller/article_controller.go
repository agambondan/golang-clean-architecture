package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
	"go-blog-api/app/utils"
	"net/http"
	"strconv"
)

type articleController struct {
	articleService service.ArticleService
	userService    service.UserService
	redis          security.Interface
	auth           security.TokenInterface
}

type ArticleController interface {
	SaveArticle(c *gin.Context)
	GetArticles(c *gin.Context)
	GetArticle(c *gin.Context)
	GetArticleByTitle(c *gin.Context)
	GetArticlesByUserId(c *gin.Context)
	GetArticlesByUsername(c *gin.Context)
	GetArticlesByCategoryName(c *gin.Context)
	GetCountArticlesByCategoryName(c *gin.Context)
	UpdateArticle(c *gin.Context)
	DeleteArticle(c *gin.Context)
	CountArticles(c *gin.Context)
}

func NewArticleController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) ArticleController {
	newArticleService := service.NewArticleService(repo.Article)
	newUserService := service.NewUserService(repo.User)
	return &articleController{newArticleService, newUserService, redis, auth}
}

func (p *articleController) SaveArticle(ctx *gin.Context) {
	var article model.Article
	var articleAPI model.ArticleAPI
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	categoriesFormArray := ctx.PostFormArray("categories")
	var categories = make([]model.Category, len(categoriesFormArray))
	for i, value := range categoriesFormArray {
		categoryIDString, _ := strconv.Atoi(value)
		categoryID := int64(categoryIDString)
		categories[i].ID = &categoryID
	}
	article.Categories = &categories
	contentType := ctx.ContentType()
	if contentType != "application/json" {
		data := ctx.PostForm("data")
		err = json.Unmarshal([]byte(data), &articleAPI)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	} else {
		if err = ctx.ShouldBindJSON(&article); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	}
	_ = lib.Merge(articleAPI, &article)
	validate := article.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
		return
	}
	fileHeader, _ := ctx.FormFile("images")
	if contentType != "application/json" && fileHeader.Filename != "" {
		uploadFile, err := utils.UploadImageFileToAssets(ctx, "article", checkIdUser.ID.String(), utils.DriveImagesId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed upload image to google drive", "user failed to created", article))
			return
		}
		article.Image = &uploadFile.Name
		imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
		article.ImageURL = &imageURL
		article.ThumbnailURL = &uploadFile.ThumbnailLink
		validate = article.Validate("images")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
			return
		}
	}
	article.UserID = checkIdUser.ID
	_, err = p.articleService.Create(&article)
	if err != nil {
		ctx.JSON(http.StatusConflict, model.BuildErrorResponse("failed create article", err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article))
}

func (p *articleController) GetArticles(ctx *gin.Context) {
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("articles not found", err.Error(), nil))
		return
	}
	var articlesType *model.Articles
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	_ = lib.Merge(articles, &articlesType)
	if *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articlesType.PublicArticles()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articles))
	}
}

func (p *articleController) GetArticle(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	article, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("article not found", err.Error(), nil))
		return
	}
	userFindById, err := p.userService.FindById(article.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("user who created article not found", err.Error(), nil))
		return
	}
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	article.Author = userFindById
	if *checkIdUser.RoleId != 1 || article.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article.PublicArticle()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article))
	}
}

func (p *articleController) GetArticleByTitle(ctx *gin.Context) {
	title := ctx.Param("title")
	article, err := p.articleService.FindByTitle(title)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("article not found", err.Error(), nil))
		return
	}
	_, err = p.userService.FindById(article.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("article not found", err.Error(), nil))
		return
	}
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if *checkIdUser.RoleId != 1 || article.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article.PublicArticle()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article))
	}
}

func (p *articleController) GetArticlesByUserId(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id := uuid.MustParse(idParam)
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAllByUserId(id, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("articles not found", err.Error(), nil))
		return
	}
	var articlesType *model.Articles
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	_ = lib.Merge(articles, &articlesType)
	if *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articlesType.PublicArticles()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articles))
	}
}

func (p *articleController) GetArticlesByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAllByUsername(username, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("articles not found", err.Error(), nil))
		return
	}
	var articlesType *model.Articles
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	_ = lib.Merge(articles, &articlesType)
	if *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articlesType.PublicArticles()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articles))
	}
}

func (p *articleController) GetArticlesByCategoryName(ctx *gin.Context) {
	limit, offset := utils.GetLimitOffsetParam(ctx)
	name := ctx.Param("name")
	articles, err := p.articleService.FindAllByCategoryName(name, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("articles not found", err.Error(), nil))
		return
	}
	var articlesType *model.Articles
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	_ = lib.Merge(articles, &articlesType)
	if *checkIdUser.RoleId != 1 {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articlesType.PublicArticles()))
	} else {
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", articles))
	}
}

func (p *articleController) GetCountArticlesByCategoryName(ctx *gin.Context) {
	name := ctx.Param("name")
	count, err := p.articleService.CountByCategoryName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, count)

}

func (p *articleController) UpdateArticle(ctx *gin.Context) {
	var article model.Article
	var articleAPI model.ArticleAPI
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	articleFindById, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.BuildErrorResponse("article not found", err.Error(), nil))
		return
	}
	if articleFindById.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", "user id not equals with article", nil))
		return
	}
	contentType := ctx.ContentType()
	if contentType != "application/json" {
		data := ctx.PostForm("data")
		if err = json.Unmarshal([]byte(data), &articleAPI); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	} else {
		if err = ctx.ShouldBindJSON(&article); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
			return
		}
	}
	_ = lib.Merge(articleAPI, &article)
	validate := article.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("fill your empty field", "field can't empty", validate))
		return
	}
	_, err = p.articleService.UpdateById(int64(id), &article)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("can't update article", err.Error(), article))
		return
	}
	fileHeader, _ := ctx.FormFile("images")
	if contentType != "application/json" && fileHeader.Filename != "" {
		uploadFile, err := utils.UploadImageFileToAssets(ctx, "article", checkIdUser.ID.String(), utils.DriveImagesId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.BuildErrorResponse("failed upload image to google drive", "user failed to created", article))
			return
		}
		article.CreatedAt = articleFindById.CreatedAt
		article.Image = &uploadFile.Name
		imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
		article.ImageURL = &imageURL
		article.ThumbnailURL = &uploadFile.ThumbnailLink
		_, err = p.articleService.UpdateById(int64(id), &article)
	}
	ctx.JSON(http.StatusOK, model.BuildResponse(true, "success", article))

}

func (p *articleController) DeleteArticle(ctx *gin.Context) {
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("id must number", err.Error(), nil))
		return
	}
	articleFindById, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("article not found", err.Error(), nil))
		return
	}
	if articleFindById.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", "user id not equals with article", nil))
		return
	} else {
		err = p.articleService.DeleteById(int64(id))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.BuildErrorResponse("article not found", err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusOK, model.BuildResponse(true, "successfully delete article", nil))
		return
	}
}

func (p *articleController) CountArticles(ctx *gin.Context) {
	count, err := p.articleService.Count()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, count)
}
