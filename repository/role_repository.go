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
	queryInsert := fmt.Sprintf("INSERT INTO %s (id, name, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5)", "roles")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return role, err
	}
	_, err = stmt.Exec(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt, nil)
	if err != nil {
		return role, err
	}
	return role, err
}

func (r *roleRepo) FindAll() ([]model.Role, error) {
	panic("implement me")
}

func (r *roleRepo) FindById(id uint64) (model.Role, error) {
	var role model.Role
	querySelect := fmt.Sprint("SELECT id, name, created_at, updated_at FROM roles WHERE id=$1 AND deleted_at IS NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return role, err
	}
	err = prepare.QueryRow(id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return role, err
	}
	return role, nil
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
