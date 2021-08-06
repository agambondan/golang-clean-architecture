package service

import (
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
)

type categoryService struct {
	category repository.CategoryRepository
}

// NewCategoryService implements the CategoryService Interface
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

type CategoryService interface {
	Create(category *model.Category) (*model.Category, error)
	FindAll() ([]model.Category, error)
	FindById(id uint64) (model.Category, error)
	FindByName(name string) (model.Category, error)
	UpdateById(id uint64, category *model.Category) (*model.Category, error)
	DeleteById(id uint64) error
}

func (p *categoryService) Create(category *model.Category) (*model.Category, error) {
	return p.category.Save(category)
}

func (p *categoryService) FindAll() ([]model.Category, error) {
	return p.category.FindAll()
}

func (p *categoryService) FindById(id uint64) (model.Category, error) {
	return p.category.FindById(id)
}

func (p *categoryService) FindByName(name string) (model.Category, error) {
	return p.category.FindByName(name)
}

func (p *categoryService) UpdateById(id uint64, category *model.Category) (*model.Category, error) {
	return p.category.UpdateById(id, category)
}

func (p *categoryService) DeleteById(id uint64) error {
	return p.category.DeleteById(id)
}
