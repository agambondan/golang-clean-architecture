package repository

import (
	"database/sql"
	"fmt"
	"golang-youtube-api/model"
)

type RoleRepository interface {
	Save(role *model.Role) (*model.Role, error)
	FindAll() ([]model.Role, error)
	FindById(id uint64) (model.Role, error)
	UpdateById(id uint64, role *model.Role) (*model.Role, error)
	DeleteById(id uint64) error
}

type roleRepo struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepo{db}
}

func (r *roleRepo) Save(role *model.Role) (*model.Role, error) {
	panic("implement me")
}

func (r *roleRepo) FindAll() ([]model.Role, error) {
	panic("implement me")
}

func (r *roleRepo) FindById(id uint64) (model.Role, error) {
	panic("implement me")
}

func (r *roleRepo) UpdateById(id uint64, role *model.Role) (*model.Role, error) {
	panic("implement me")
}

func (r *roleRepo) DeleteById(id uint64) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where id = %d", "roles", id)
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
