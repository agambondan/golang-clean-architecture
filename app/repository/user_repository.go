package repository

import (
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *model.User) (*model.User, error)
	FindAll(limit, offset int) (*[]model.User, error)
	FindById(uuid *uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAllByRoleId(id int64) (*[]model.User, error)
	FindUserByEmailAndPassword(user *model.User) (*model.User, error)
	UpdateById(uuid *uuid.UUID, user *model.User) (*model.User, error)
	DeleteById(uuid *uuid.UUID) error
	Count() (int64, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (u *userRepo) Save(user *model.User) (*model.User, error) {
	if tx := u.db.Create(&user); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindAll(limit, offset int) (*[]model.User, error) {
	var users *[]model.User
	if tx := u.db.
		//Joins("join article a on \"user\".id = a.user_id").Joins("join role r on \"user\".role_id = r.id").
		Preload("Role").
		Find(&users).Offset(offset).Limit(limit); tx.Error != nil {
		return users, tx.Error
	}
	return users, nil
}

func (u *userRepo) FindById(uuid *uuid.UUID) (*model.User, error) {
	var user *model.User
	if tx := u.db.
		//Joins("join article a on \"user\".id = a.user_id").Joins("join role r on \"user\".role_id = r.id").
		First(&user, uuid); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindByUsername(username string) (*model.User, error) {
	var user *model.User
	if tx := u.db.
		//Joins("join article a on user.id = a.user_id").Joins("join role r on user.role_id = r.id").
		First(&user, "user.username = ?", username); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindAllByRoleId(id int64) (*[]model.User, error) {
	var users *[]model.User
	if tx := u.db.
		//Joins("join article a on user.id = a.user_id").Joins("join role r on user.role_id = r.id").
		Find(&users, "user.role_id = ?", id); tx.Error != nil {
		return users, tx.Error
	}
	return users, nil
}

func (u *userRepo) FindUserByEmailAndPassword(userToken *model.User) (*model.User, error) {
	var user *model.User
	if tx := u.db.First(&user, "(username = ? or email = ?) and password = ?", userToken.Username, userToken.Email, userToken.Password); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) UpdateById(uuid *uuid.UUID, user *model.User) (*model.User, error) {
	findById, err := u.FindById(uuid)
	if err != nil {
		return findById, err
	}
	_ = lib.Merge(findById, &user)
	if tx := u.db.Updates(&user); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) DeleteById(uuid *uuid.UUID) error {
	_, err := u.FindById(uuid)
	if err != nil {
		return err
	}
	u.db.Delete(&model.User{}, uuid)
	return nil
}

func (u *userRepo) Count() (int64, error) {
	var count int64
	u.db.Model(&[]model.User{}).Count(&count)
	return count, nil
}
