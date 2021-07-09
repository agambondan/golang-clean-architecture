package repository

import (
	"database/sql"
	"fmt"
	"golang-youtube-api/model"
)

type CategoryRepository interface {
	Save(category *model.Category) (*model.Category, error)
	FindAll() ([]model.Category, error)
	FindById(id uint64) (model.Category, error)
	UpdateById(id uint64, category *model.Category) (*model.Category, error)
	DeleteById(id uint64) error
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{db}
}

func (r *categoryRepo) Save(category *model.Category) (*model.Category, error) {
	panic("implement me")
}

func (r *categoryRepo) FindAll() ([]model.Category, error) {
	panic("implement me")
}

func (r *categoryRepo) FindById(id uint64) (model.Category, error) {
	panic("implement me")
}

func (r *categoryRepo) UpdateById(id uint64, category *model.Category) (*model.Category, error) {
	panic("implement me")
}

func (r *categoryRepo) DeleteById(id uint64) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where id = %d", "category", id)
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
