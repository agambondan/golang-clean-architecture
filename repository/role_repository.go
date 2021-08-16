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
	Count() (int, error)
}

type roleRepo struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepo{db}
}

func (r *roleRepo) Save(role *model.Role) (*model.Role, error) {
	queryInsert := fmt.Sprintf("insert into %s (name, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4) returning id", "roles")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return role, err
	}
	err = stmt.QueryRow(&role.Name, &role.CreatedAt, &role.UpdatedAt, nil).Scan(&role.ID)
	if err != nil {
		return role, err
	}
	return role, err
}

func (r *roleRepo) FindAll() ([]model.Role, error) {
	var roles []model.Role
	var role model.Role
	query := fmt.Sprintf("select id, name, created_at, updated_at from roles where deleted_at is null")
	rows, err := r.db.Query(query)
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		err = rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, err
}

func (r *roleRepo) FindById(id uint64) (model.Role, error) {
	var role model.Role
	querySelect := fmt.Sprint("select id, name, created_at, updated_at from roles where id=$1 and deleted_at is null")
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
	query := fmt.Sprintf("update roles set name = $1, updated_at = $2 where id = %d", id)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return role, err
	}
	_, err = stmt.Exec(&role.Name, &role.UpdatedAt)
	role.ID = id
	if err != nil {
		return role, err
	}
	return role, err
}

func (r *roleRepo) DeleteById(id uint64) error {
	query := fmt.Sprintf("delete from roles where id = %d", id)
	_, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	r.db.Exec(query)
	return err
}

func (r *roleRepo) Count() (int, error) {
	var count int
	queryInsert := fmt.Sprintf("select count(id) from roles")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return count, err
	}
	err = prepare.QueryRow().Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}
