package service

import (
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
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
	FindAll(limit, offset int) (*[]model.Category, error)
	FindAllByWord(word string, limit, offset int) (*[]model.Category, error)
	FindById(id int64) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	UpdateById(id int64, category *model.Category) (*model.Category, error)
	DeleteById(id int64) error
	Count() (int64, error)
}

func (p *categoryService) Create(category *model.Category) (*model.Category, error) {
	return p.category.Save(category)
}

func (p *categoryService) FindAll(limit, offset int) (*[]model.Category, error) {
	return p.category.FindAll(limit, offset)
}

func (p *categoryService) FindAllByWord(word string, limit, offset int) (*[]model.Category, error) {
	return p.category.FindAllByWord(word, limit, offset)
}

func (p *categoryService) FindById(id int64) (*model.Category, error) {
	return p.category.FindById(id)
}

func (p *categoryService) FindByName(name string) (*model.Category, error) {
	return p.category.FindByName(name)
}

func (p *categoryService) UpdateById(id int64, category *model.Category) (*model.Category, error) {
	return p.category.UpdateById(id, category)
}

func (p *categoryService) DeleteById(id int64) error {
	return p.category.DeleteById(id)
}

func (p *categoryService) Count() (int64, error) {
	return p.category.Count()
}
