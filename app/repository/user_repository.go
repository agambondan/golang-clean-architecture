package repository

import (
	"fmt"
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"gorm.io/gorm"
	"os"
)

type UserRepository interface {
	Save(user *model.User) (*model.User, error)
	FindAll(limit, offset int) (*[]model.User, error)
	FindById(uuid *uuid.UUID) (*model.User, error)
	FindAllByUsername(username string) (*[]model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAllByRoleId(id int64, offset, limit int) (*[]model.User, error)
	FindUserByEmailOrUsername(user *model.User) (*model.User, error)
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
	cipherEncrypt, err := lib.CipherEncrypt([]byte(*user.Password), []byte(os.Getenv("CIPHER_KEY")))
	if err != nil {
		return user, err
	}
	cipherEncryptString := fmt.Sprintf("%x", cipherEncrypt)
	user.Password = &cipherEncryptString
	if tx := u.db.Create(&user); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindAll(limit, offset int) (*[]model.User, error) {
	var users *[]model.User
	if tx := u.db.Preload("Role").Limit(limit).Offset(offset).Find(&users); tx.Error != nil {
		return users, tx.Error
	}
	return users, nil
}

func (u *userRepo) FindById(uuid *uuid.UUID) (*model.User, error) {
	var user *model.User
	if tx := u.db.Preload("Role").First(&user, uuid); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindAllByUsername(username string) (*[]model.User, error) {
	var user *[]model.User
	if tx := u.db.Preload("Role").Find(&user, fmt.Sprint("username like '%"+username+"%'")); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindByUsername(username string) (*model.User, error) {
	var user *model.User
	if tx := u.db.Preload("Role").First(&user, "username = ?", username); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) FindAllByRoleId(id int64, offset, limit int) (*[]model.User, error) {
	var users *[]model.User
	if tx := u.db.Preload("Role").Limit(limit).Offset(offset).Find(&users, "role_id = ?", id); tx.Error != nil {
		return users, tx.Error
	}
	return users, nil
}

func (u *userRepo) FindUserByEmailOrUsername(userToken *model.User) (*model.User, error) {
	var user *model.User
	if tx := u.db.First(&user, "username = ? or email = ?", userToken.Username, userToken.Email); tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (u *userRepo) UpdateById(uuid *uuid.UUID, user *model.User) (*model.User, error) {
	cipherEncrypt, err := lib.CipherEncrypt([]byte(*user.Password), []byte(os.Getenv("CIPHER_KEY")))
	if err != nil {
		return user, err
	}
	cipherEncryptString := fmt.Sprintf("%x", cipherEncrypt)
	user.Password = &cipherEncryptString
	_, err = u.FindById(uuid)
	if err != nil {
		return user, err
	}
	user.ID = uuid
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
	u.db.Table("user").Select("id").Count(&count)
	return count, nil
}
