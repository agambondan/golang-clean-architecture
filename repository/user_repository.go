package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang-youtube-api/model"
)

type UserRepository interface {
	Save(user *model.User) (*model.User, error)
	FindAll(limit, offset int) ([]model.User, error)
	FindById(uuid uuid.UUID) (model.User, error)
	FindByUsername(username string) (model.User, error)
	FindAllByRoleId(id uint64) ([]model.User, error)
	FindUserByEmailAndPassword(user *model.User) (model.User, error)
	UpdateById(uuid uuid.UUID, user *model.User) (*model.User, error)
	DeleteById(uuid uuid.UUID) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Save(user *model.User) (*model.User, error) {
	queryInsert := fmt.Sprintf("INSERT INTO %s (uuid, first_name, last_name, email, phone_number, username, password, photo_profile, role_id, "+
		"instagram, facebook, twitter, linkedin, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", "users")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return user, err
	}
	_, err = stmt.Exec(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile, &user.RoleId,
		&user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt, nil)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *userRepo) FindAll(limit, offset int) ([]model.User, error) {
	var users []model.User
	queryGetUsers := fmt.Sprintf("SELECT uuid, first_name, last_name, email, phone_number, username, password, "+
		"photo_profile, role_id, instagram, facebook, twitter, linkedin, created_at, updated_at FROM users WHERE deleted_at IS NULL "+
		"limit %d offset %d ;", limit, offset)
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile, &user.RoleId,
			&user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) FindById(uuid uuid.UUID) (model.User, error) {
	var user model.User
	querySelect := fmt.Sprint("SELECT uuid, first_name, last_name, email, phone_number, username, password, photo_profile, " +
		"role_id, instagram, facebook, twitter, linkedin, created_at, updated_at FROM users WHERE uuid=$1 AND deleted_at IS NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return user, err
	}
	err = prepare.QueryRow(uuid).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile,
		&user.RoleId, &user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) FindByUsername(username string) (model.User, error) {
	var user model.User
	querySelect := fmt.Sprint("SELECT uuid, first_name, last_name, email, phone_number, username, password, photo_profile, " +
		"role_id, instagram, facebook, twitter, linkedin, created_at, updated_at FROM users WHERE username=$1 AND deleted_at IS NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return user, err
	}
	err = prepare.QueryRow(&username).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile,
		&user.RoleId, &user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) FindAllByRoleId(id uint64) ([]model.User, error) {
	var users []model.User
	query := fmt.Sprintf("select u.uuid, u.first_name, u.last_name, u.email, u.phone_number, u.username, u.photo_profile, u.password, u.role_id, "+
		"u.instagram, u.facebook, u.twitter, u.linkedin, u.created_at, u.updated_at from users u inner join roles r on r.id = u.role_id where r.id = %d and u.deleted_at is null", id)
	rows, err := r.db.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile, &user.RoleId,
			&user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) FindUserByEmailAndPassword(u *model.User) (model.User, error) {
	var user model.User
	queryLogin := fmt.Sprint("SELECT uuid, first_name, last_name, email, phone_number, username, password, photo_profile, role_id, instagram, facebook, twitter, linkedin, " +
		"created_at, updated_at FROM users WHERE (email=$1 OR username=$2) AND password=$3 AND deleted_at IS NULL")
	err := r.db.QueryRow(queryLogin, u.Email, u.Username, u.Password).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username,
		&user.Password, &user.PhotoProfile, &user.RoleId, &user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) UpdateById(uuid uuid.UUID, user *model.User) (*model.User, error) {
	queryInsert := fmt.Sprint("UPDATE users SET first_name = $1, last_name = $2, email = $3, phone_number = $4," +
		"username = $5, password = $6, image = $7, role_id = $8, updated_at = $9 where uuid = $10")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return user, err
	}
	_, err = stmt.Exec(&user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.PhotoProfile, &user.RoleId,
		&user.Instagram, &user.Facebook, &user.Twitter, &user.LinkedIn, &user.UpdatedAt, uuid.String())
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *userRepo) DeleteById(uuid uuid.UUID) error {
	queryInsert := fmt.Sprintf("DELETE FROM %s where uuid = %s", "users", uuid.String())
	_, err := r.db.Prepare(queryInsert)
	if err != nil {
		return err
	}
	return err
}
