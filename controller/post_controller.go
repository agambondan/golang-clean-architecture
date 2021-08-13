package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"golang-youtube-api/utils"
	"golang-youtube-api/utils/google"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type postController struct {
	postService         service.PostService
	userService         service.UserService
	postCategoryService service.PostCategoryService
	redis               security.Interface
	auth                security.TokenInterface
}

type PostController interface {
	SavePost(c *gin.Context)
	GetPosts(c *gin.Context)
	GetPost(c *gin.Context)
	GetPostByTitle(c *gin.Context)
	GetPostsByUserId(c *gin.Context)
	GetPostsByUsername(c *gin.Context)
	GetPostsByCategoryName(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

func NewPostController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) PostController {
	newPostService := service.NewPostService(repo.Post)
	newUserService := service.NewUserService(repo.User)
	newPostCategoryService := service.NewPostCategoryService(repo.PostCategory)
	return &postController{newPostService, newUserService, newPostCategoryService, redis, auth}
}

func (p *postController) SavePost(ctx *gin.Context) {
	var post model.Post
	post.Prepare()
	checkIdUser, err := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	contentType := ctx.ContentType()
	categoryArray := ctx.PostFormArray("categories")
	if contentType != "application/json" {
		post.Title = ctx.PostForm("title")
		post.Description = ctx.PostForm("description")
	} else {
		if err = ctx.ShouldBindJSON(&post); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "invalid json",
			})
			return
		}
	}
	post.UserUUID = checkIdUser.UUID
	post.Author = checkIdUser
	uploadFile, err := google.UploadImageFileToAssets(ctx, "post", checkIdUser.UUID.String(), utils.DriveImagesId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	post.Image = uploadFile.Name
	post.ImageURL = fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", uploadFile.Id)
	post.ThumbnailURL = uploadFile.ThumbnailLink
	validate := post.Validate("")
	if len(validate) != 0 {
		ctx.JSON(http.StatusBadRequest, validate)
		return
	}
	postCreate, err := p.postService.Create(&post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	for j := 0; j < len(categoryArray); j++ {
		var postCategory = model.PostCategory{}
		postCategory.CategoryID, _ = strconv.ParseUint(categoryArray[j], 10, 64)
		postCategory.PostID = post.ID
		_, err = p.postCategoryService.Create(&postCategory)
		if err != nil {
			log.Println(err)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": postCreate})
}

func (p *postController) GetPosts(ctx *gin.Context) {
	posts := model.Posts{}
	var limit, offset int
	var err error
	queryParamLimit := ctx.Query("_limit")
	queryParamOffset := ctx.Query("_offset")
	if queryParamLimit != "" {
		limit, err = strconv.Atoi(queryParamLimit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	} else {
		limit = 4
	}
	if queryParamOffset != "" {
		offset, err = strconv.Atoi(queryParamOffset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	}
	posts, err = p.postService.FindAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, posts.PublicPosts())
}

func (p *postController) GetPost(ctx *gin.Context) {
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	post, err := p.postService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userFindById, err := p.userService.FindById(post.UserUUID)
	if err != nil {
		return
	}
	post.Author = userFindById
	if checkIdUser.RoleId != 1 || post.UserUUID != checkIdUser.UUID {
		ctx.JSON(http.StatusOK, post.PublicPost())
		return
	} else {
		ctx.JSON(http.StatusOK, post)
	}
}

func (p *postController) GetPostByTitle(ctx *gin.Context) {
	checkIdUser, _ := utils.CheckIdUser(p.auth, p.redis, p.userService, ctx)
	title := ctx.Param("title")
	post, err := p.postService.FindByTitle(title)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userFindById, err := p.userService.FindById(post.UserUUID)
	if err != nil {
		return
	}
	post.Author = userFindById
	if checkIdUser.RoleId != 1 || post.UserUUID != checkIdUser.UUID {
		ctx.JSON(http.StatusOK, post.PublicPost())
		return
	} else {
		ctx.JSON(http.StatusOK, post)
	}
}

func (p *postController) GetPostsByUserId(ctx *gin.Context) {
	posts := model.Posts{}
	idParam := ctx.Param("id")
	id := uuid.MustParse(idParam)
	posts, err := p.postService.FindAllByUserId(id)
	for i := 0; i < len(posts); i++ {
		posts[i].Author, err = p.userService.FindById(posts[i].UserUUID)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func (p *postController) GetPostsByUsername(ctx *gin.Context) {
	posts := model.Posts{}
	username := ctx.Param("username")
	posts, err := p.postService.FindAllByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, posts.PublicPosts())
}

func (p *postController) GetPostsByCategoryName(ctx *gin.Context) {
	posts := model.Posts{}
	var limit, offset int
	var err error
	queryParamLimit := ctx.Query("_limit")
	queryParamOffset := ctx.Query("_offset")
	if queryParamLimit != "" {
		limit, err = strconv.Atoi(queryParamLimit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	} else {
		limit = 4
	}
	if queryParamOffset != "" {
		offset, err = strconv.Atoi(queryParamOffset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
	}
	name := ctx.Param("name")
	posts, err = p.postService.FindAllByCategoryName(name, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, posts.PublicPosts())
}

func (p *postController) UpdatePost(ctx *gin.Context) {
	var post model.Post
	post.Prepare()
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
	postFindById, err := p.postService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if postFindById.UserUUID != checkIdUser.UUID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	} else {
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		categoryArray := ctx.PostFormArray("categories")
		if title != "" && description != "" {
			post.Title = title
			post.Description = description
		} else {
			if err = ctx.ShouldBindJSON(&post); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "invalid json",
				})
				return
			}
		}
		post.UserUUID = checkIdUser.UUID
		validate := post.Validate("")
		if len(validate) != 0 {
			ctx.JSON(http.StatusBadRequest, validate)
			return
		}
		post.Author = checkIdUser
		postUpdate, err := p.postService.UpdateById(uint64(id), &post)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		filenames, err := utils.CreateUploadPhotoMachine(ctx, post.UserUUID.String(), "/post")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		postUpdate.CreatedAt = postFindById.CreatedAt
		post.Image = strings.Join(filenames, "")
		for j := 0; j < len(categoryArray); j++ {
			var postCategory = model.PostCategory{}
			postCategory.CategoryID, _ = strconv.ParseUint(categoryArray[j], 10, 64)
			postCategory.PostID = post.ID
			_, err = p.postCategoryService.UpdateById(uint64(id), &postCategory)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"data": postUpdate, "filenames": filenames})
	}
}

func (p *postController) DeletePost(ctx *gin.Context) {
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
	postFindById, err := p.postService.FindById(uint64(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if postFindById.UserUUID != checkIdUser.UUID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	} else {
		err = p.postService.DeleteById(uint64(id))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete Successfully"})
		return
	}
}
