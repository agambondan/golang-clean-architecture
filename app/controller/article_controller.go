package controller

import (
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
	"strings"
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
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		article.Title = &title
		article.Description = &description
	} else {
		if err = ctx.ShouldBindJSON(&article); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "invalid json",
			})
			return
		}
	}
	validate := article.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, validate)
		return
	}
	article.UserID = checkIdUser.ID
	article.Author = checkIdUser
	uploadFile, err := utils.UploadImageFileToAssets(ctx, "article", checkIdUser.ID.String(), utils.DriveImagesId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	article.Image = &uploadFile.Name
	imageURL := fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
	article.ImageURL = &imageURL
	article.ThumbnailURL = &uploadFile.ThumbnailLink
	article.Validate("images")
	articleCreate, err := p.articleService.Create(&article)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": articleCreate})
}

func (p *articleController) GetArticles(ctx *gin.Context) {
	var articlesType *model.Articles
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	fmt.Println(checkIdUser, "JANCOK")
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	_ = lib.Merge(articles, &articlesType)
	if checkIdUser != nil {
		if *checkIdUser.RoleId != 1 {
			ctx.JSON(http.StatusOK, articlesType.PublicArticles())
		} else {
			ctx.JSON(http.StatusOK, articles)
		}
	} else {
		fmt.Println("JANCOK")
		ctx.JSON(http.StatusOK, articlesType.PublicArticles())
	}
}

func (p *articleController) GetArticle(ctx *gin.Context) {
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	article, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userFindById, err := p.userService.FindById(article.UserID)
	if err != nil {
		return
	}
	article.Author = userFindById
	if *checkIdUser.RoleId != 1 || article.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusOK, article.PublicArticle())
		return
	} else {
		ctx.JSON(http.StatusOK, article)
	}
}

func (p *articleController) GetArticleByTitle(ctx *gin.Context) {
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	title := ctx.Param("title")
	article, err := p.articleService.FindByTitle(title)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err = p.userService.FindById(article.UserID)
	if err != nil {
		return
	}
	if *checkIdUser.RoleId != 1 || article.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusOK, article.PublicArticle())
	} else {
		ctx.JSON(http.StatusOK, article)
	}
}

func (p *articleController) GetArticlesByUserId(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id := uuid.MustParse(idParam)
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAllByUserId(id, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func (p *articleController) GetArticlesByUsername(ctx *gin.Context) {
	articlesType := model.Articles{}
	username := ctx.Param("username")
	limit, offset := utils.GetLimitOffsetParam(ctx)
	articles, err := p.articleService.FindAllByUsername(username, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	_ = lib.Merge(articles, &articlesType)
	ctx.JSON(http.StatusOK, articlesType.PublicArticles())
}

func (p *articleController) GetArticlesByCategoryName(ctx *gin.Context) {
	articlesType := model.Articles{}
	limit, offset := utils.GetLimitOffsetParam(ctx)
	name := ctx.Param("name")
	articles, err := p.articleService.FindAllByCategoryName(name, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	_ = lib.Merge(articles, &articlesType)
	ctx.JSON(http.StatusOK, articlesType.PublicArticles())
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
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	articleFindById, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if articleFindById.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	} else {
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		if title != "" && description != "" {
			article.Title = &title
			article.Description = &description
		} else {
			if err = ctx.ShouldBindJSON(&article); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid json",
				})
				return
			}
		}
		article.UserID = checkIdUser.ID
		validate := article.Validate("")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, validate)
			return
		}
		//article.Author = &checkIdUser
		articleUpdate, err := p.articleService.UpdateById(int64(id), &article)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		filenames, err := utils.CreateUploadPhotoMachine(ctx, article.UserID.String(), "/article")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		articleUpdate.CreatedAt = articleFindById.CreatedAt
		image := strings.Join(filenames, "")
		article.Image = &image
		ctx.JSON(http.StatusOK, gin.H{"data": articleUpdate, "filenames": filenames})
	}
}

func (p *articleController) DeleteArticle(ctx *gin.Context) {
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	articleFindById, err := p.articleService.FindById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if articleFindById.UserID != checkIdUser.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	} else {
		err = p.articleService.DeleteById(int64(id))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete Successfully"})
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
