package service

import (
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
)

type postService struct {
	post repository.PostRepository
}

// NewPostService implements the PostService Interface
func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo}
}

type PostService interface {
	Create(post *model.Post) (*model.Post, error)
	FindAll(limit, offset int) ([]model.Post, error)
	FindById(id uint64) (model.Post, error)
	FindByTitle(title string) (model.Post, error)
	FindAllByUserId(uuid uuid.UUID) ([]model.Post, error)
	FindAllByUsername(username string) ([]model.Post, error)
	FindAllByCategoryName(name string, limit, offset int) ([]model.Post, error)
	UpdateById(id uint64, post *model.Post) (*model.Post, error)
	DeleteById(id uint64) error
	Count() (int, error)
}

func (p *postService) Create(post *model.Post) (*model.Post, error) {
	return p.post.Save(post)
}

func (p *postService) FindAll(limit, offset int) ([]model.Post, error) {
	return p.post.FindAll(limit, offset)
}

func (p *postService) FindById(id uint64) (model.Post, error) {
	return p.post.FindById(id)
}

func (p *postService) FindByTitle(title string) (model.Post, error) {
	return p.post.FindByTitle(title)
}

func (p *postService) FindAllByUserId(uuid uuid.UUID) ([]model.Post, error) {
	return p.post.FindAllByUserId(uuid)
}

func (p *postService) FindAllByUsername(username string) ([]model.Post, error) {
	return p.post.FindAllByUsername(username)
}

func (p *postService) FindAllByCategoryName(name string, limit, offset int) ([]model.Post, error) {
	return p.post.FindAllByCategoryName(name, limit, offset)
}

func (p *postService) UpdateById(id uint64, post *model.Post) (*model.Post, error) {
	return p.post.UpdateById(id, post)
}

func (p *postService) DeleteById(id uint64) error {
	return p.post.DeleteById(id)
}

func (p *postService) Count() (int,error) {
	return p.post.Count()
}