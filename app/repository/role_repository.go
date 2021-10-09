package repository

import (
	"go-blog-api/app/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(role *model.Role) (*model.Role, error)
	FindAll() (*[]model.Role, error)
	FindById(id *int64) (*model.Role, error)
	UpdateById(id *int64, role *model.Role) (*model.Role, error)
	DeleteById(id *int64) error
	Count() (*int64, error)
}

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepo{db}
}

func (r *roleRepo) Save(role *model.Role) (*model.Role, error) {
	if tx := r.db.Create(&role); tx.Error != nil {
		return role, tx.Error
	}
	return role, nil
}

func (r *roleRepo) FindAll() (*[]model.Role, error) {
	var roles *[]model.Role
	r.db.Model(&model.Role{}).Find(&roles)
	return roles, nil
}

func (r *roleRepo) FindById(id *int64) (*model.Role, error) {
	var role *model.Role
	if tx := r.db.First(&role, id); tx.Error != nil || tx.RowsAffected < 1 {
		return role, tx.Error
	}
	return role, nil
}

func (r *roleRepo) UpdateById(id *int64, role *model.Role) (*model.Role, error) {
	findById, err := r.FindById(id)
	if err != nil {
		return findById, err
	}
	if tx := r.db.Updates(&role); tx.Error != nil {
		return role, tx.Error
	}
	return role, nil
}

func (r *roleRepo) DeleteById(id *int64) error {
	_, err := r.FindById(id)
	if err != nil {
		return err
	}
	r.db.Delete(&model.Role{}, id)
	return nil
}

func (r *roleRepo) Count() (*int64, error) {
	var count int64
	r.db.Table("role").Select("id").Count(&count)
	return &count, nil
}
