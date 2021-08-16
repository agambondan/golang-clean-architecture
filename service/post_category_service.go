package service

import (
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
)

type postCategoryService struct {
	post repository.PostCategoryRepository
}

// NewPostCategoryService implements the PostCategoryService Interface
func NewPostCategoryService(repo repository.PostCategoryRepository) PostCategoryService {
	return &postCategoryService{repo}
}

type PostCategoryService interface {
	Create(post *model.PostCategory) (*model.PostCategory, error)
	FindAll() ([]model.PostCategory, error)
	UpdateById(id uint64, post *model.PostCategory) (*model.PostCategory, error)
	DeleteById(id uint64) error
	Count() (int, error)
}

func (p *postCategoryService) Create(post *model.PostCategory) (*model.PostCategory, error) {
	return p.post.Save(post)
}

func (p *postCategoryService) FindAll() ([]model.PostCategory, error) {
	return p.post.FindAll()
}

func (p *postCategoryService) UpdateById(id uint64, post *model.PostCategory) (*model.PostCategory, error) {
	return p.post.UpdateById(id, post)
}

func (p *postCategoryService) DeleteById(id uint64) error {
	return p.post.DeleteById(id)
}

func (p *postCategoryService) Count() (int,error) {
	return p.post.Count()
}