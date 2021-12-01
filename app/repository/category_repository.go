package repository

import (
	"go-blog-api/app/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Save(category *model.Category) (*model.Category, error)
	FindAll(limit, offset int) (*[]model.Category, error)
	FindAllByWord(word string, limit, offset int) (*[]model.Category, error)
	FindById(id int64) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	UpdateById(id int64, category *model.Category) (*model.Category, error)
	DeleteById(id int64) error
	Count() (int64, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db}
}

func (c *categoryRepo) Save(category *model.Category) (*model.Category, error) {
	if tx := c.db.Create(&category); tx.Error != nil {
		return category, tx.Error
	}
	return category, nil
}

func (c *categoryRepo) FindAll(limit, offset int) (*[]model.Category, error) {
	var categories *[]model.Category
	c.db.Model(&model.Category{}).Offset(offset).Limit(limit).Find(&categories)
	return categories, nil
}

func (c *categoryRepo) FindAllByWord(word string, limit, offset int) (*[]model.Category, error) {
	var categories *[]model.Category
	c.db.Model(&model.Category{}).Where("category.name LIKE ?", "%"+word+"%").Offset(offset).Limit(limit).Find(&categories)
	return categories, nil
}

func (c *categoryRepo) FindById(id int64) (*model.Category, error) {
	var category *model.Category
	if tx := c.db.First(&category, id); tx.Error != nil || tx.RowsAffected < 1 {
		return category, tx.Error
	}
	return category, nil
}

func (c *categoryRepo) FindByName(name string) (*model.Category, error) {
	var category *model.Category
	if tx := c.db.First(&category, "name = ?", name); tx.Error != nil || tx.RowsAffected < 1 {
		return category, tx.Error
	}
	return category, nil
}

func (c *categoryRepo) UpdateById(id int64, category *model.Category) (*model.Category, error) {
	findById, err := c.FindById(id)
	if err != nil {
		return findById, err
	}
	category.ID = findById.ID
	if tx := c.db.Updates(&category); tx.Error != nil {
		return category, tx.Error
	}
	return category, nil
}

func (c *categoryRepo) DeleteById(id int64) error {
	_, err := c.FindById(id)
	if err != nil {
		return err
	}
	c.db.Delete(&model.Category{}, id)
	return nil
}

func (c *categoryRepo) Count() (int64, error) {
	var count int64
	c.db.Table("category").Select("id").Where("deleted_at is null").Count(&count)
	return count, nil
}
