package repository

import (
	"database/sql"
	"fmt"
	"golang-youtube-api/model"
)

type CategoryRepository interface {
	Save(category *model.Category) (*model.Category, error)
	FindAll(limit, offset int) ([]model.Category, error)
	FindById(id uint64) (model.Category, error)
	FindByName(name string) (model.Category, error)
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
	queryInsert := fmt.Sprintf("insert into %s (name, image, image_url, thumbnail_url, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) returning id", "categories")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return category, err
	}
	err = stmt.QueryRow(&category.Name, &category.Image, &category.ImageURL, &category.ThumbnailURL, &category.CreatedAt, &category.UpdatedAt, nil).Scan(&category.ID)
	if err != nil {
		return category, err
	}
	return category, err
}

func (r *categoryRepo) FindAll(limit, offset int) ([]model.Category, error) {
	var categories []model.Category
	var category model.Category
	query := fmt.Sprintf("select id, name, image, image_url, thumbnail_url, created_at, updated_at from categories where deleted_at is null "+
		"order by id limit %d offset %d ;", limit, offset)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return categories, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return categories, err
	}
	for rows.Next() {
		err = rows.Scan(&category.ID, &category.Name, &category.Image, &category.ImageURL, &category.ThumbnailURL, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, err
}

func (r *categoryRepo) FindById(id uint64) (model.Category, error) {
	var category model.Category
	querySelect := fmt.Sprint("select id, name, image, image_url, thumbnail_url, created_at, updated_at from categories where id=$1 and deleted_at is null")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return category, err
	}
	err = prepare.QueryRow(id).Scan(&category.ID, &category.Name, &category.Image, &category.ImageURL, &category.ThumbnailURL, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepo) FindByName(name string) (model.Category, error) {
	var category model.Category
	querySelect := fmt.Sprint("select id, name, image, image_url, thumbnail_url, created_at, updated_at from categories where name=$1 and deleted_at is null")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return category, err
	}
	err = prepare.QueryRow(name).Scan(&category.ID, &category.Name, &category.Image, &category.ImageURL, &category.ThumbnailURL, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepo) UpdateById(id uint64, category *model.Category) (*model.Category, error) {
	query := fmt.Sprintf("update categories set name = $1, image = $2, image_url = $3, thumbanil_url = $4, updated_at = $5 where id = %d", id)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return category, err
	}
	_, err = stmt.Exec(&category.Name, &category.Image, &category.ImageURL, &category.ThumbnailURL, &category.UpdatedAt)
	category.ID = id
	if err != nil {
		return category, err
	}
	return category, err
}

func (r *categoryRepo) DeleteById(id uint64) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where id = %d", "categories", id)
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
