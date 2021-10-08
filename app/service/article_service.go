package service

import (
	"github.com/google/uuid"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
)

type articleService struct {
	article repository.ArticleRepository
}

// NewArticleService implements the ArticleService Interface
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo}
}

type ArticleService interface {
	Create(article *model.Article) (*model.Article, error)
	FindAll(limit, offset int) (*[]model.Article, error)
	FindById(id int64) (*model.Article, error)
	FindByTitle(title string) (*model.Article, error)
	FindAllByUserId(uuid uuid.UUID, limit, offset int) (*[]model.Article, error)
	FindAllByUsername(username string, limit, offset int) (*[]model.Article, error)
	FindAllByCategoryName(name string, limit, offset int) (*[]model.Article, error)
	CountByCategoryName(name string) (int64, error)
	UpdateById(id int64, article *model.Article) (*model.Article, error)
	DeleteById(id int64) error
	Count() (int64, error)
}

func (a *articleService) Create(article *model.Article) (*model.Article, error) {
	return a.article.Save(article)
}

func (a *articleService) FindAll(limit, offset int) (*[]model.Article, error) {
	return a.article.FindAll(limit, offset)
}

func (a *articleService) FindById(id int64) (*model.Article, error) {
	return a.article.FindById(id)
}

func (a *articleService) FindByTitle(title string) (*model.Article, error) {
	return a.article.FindByTitle(title)
}

func (a *articleService) FindAllByUserId(uuid uuid.UUID, limit, offset int) (*[]model.Article, error) {
	return a.article.FindAllByUserId(uuid, limit, offset)
}

func (a *articleService) FindAllByUsername(username string, limit, offset int) (*[]model.Article, error) {
	return a.article.FindAllByUsername(username, limit, offset)
}

func (a *articleService) FindAllByCategoryName(name string, limit, offset int) (*[]model.Article, error) {
	return a.article.FindAllByCategoryName(name, limit, offset)
}

func (a *articleService) CountByCategoryName(name string) (int64, error) {
	return a.article.CountByCategoryName(name)
}

func (a *articleService) UpdateById(id int64, article *model.Article) (*model.Article, error) {
	return a.article.UpdateById(id, article)
}

func (a *articleService) DeleteById(id int64) error {
	return a.article.DeleteById(id)
}

func (a *articleService) Count() (int64, error) {
	return a.article.Count()
}
