package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang-youtube-api/models"
)

type UserRepository interface {
	Save(user *models.User) (*models.User, error)
	FindAll() ([]models.User, error)
	FindById(uuid uuid.UUID) (models.User, error)
	UpdateById(uuid uuid.UUID, user *models.User) (*models.User, error)
	DeleteById(uuid uuid.UUID) error
}

type repo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &repo{db}
}

func (r *repo) Save(user *models.User) (*models.User, error) {
	queryInsert := fmt.Sprintf("INSERT INTO %s (uuid, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", "users")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return user, err
	}
	_, err = stmt.Exec(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt, nil)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *repo) FindAll() ([]models.User, error) {
	var users []models.User
	queryGetUsers := fmt.Sprintf("SELECT uuid, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at FROM users WHERE deleted_at IS NULL")
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		fmt.Println("JANCOK")
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, nil
}

func (r *repo) FindById(uuid uuid.UUID) (models.User, error) {
	var user models.User
	querySelect := fmt.Sprint("SELECT uuid, first_name, last_name, email, phone_number, username, password," +
		" role_id, created_at, updated_at FROM users WHERE uuid=$1 AND deleted_at IS NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return user, err
	}
	err = prepare.QueryRow(uuid).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber,
		&user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repo) UpdateById(uuid uuid.UUID, user *models.User) (*models.User, error) {
	queryInsert := fmt.Sprintf("UPDATE %s SET uuid = $1, first_name = $2, last_name = $3, email = $4, phone_number = $5,"+
		"username = $6, password = $7, role_id = $8, updated_at = $9 where uuid = %s", "users", uuid.String())
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return user, err
	}
	_, err = stmt.Exec(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber,
		&user.Username, &user.Password, &user.RoleId, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *repo) DeleteById(uuid uuid.UUID) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where uuid = %s", "users", uuid.String())
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
