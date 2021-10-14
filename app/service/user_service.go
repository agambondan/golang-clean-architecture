package service

import (
	"github.com/google/uuid"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
)

type userService struct {
	user repository.UserRepository
}

// NewUserService implements the UserService Interface
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

type UserService interface {
	Create(user *model.User) (*model.User, error)
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

func (s *userService) Create(user *model.User) (*model.User, error) {
	return s.user.Save(user)
}

func (s *userService) FindAll(limit, offset int) (*[]model.User, error) {
	return s.user.FindAll(limit, offset)
}

func (s *userService) FindById(uuid *uuid.UUID) (*model.User, error) {
	return s.user.FindById(uuid)
}

func (s *userService) FindAllByUsername(username string) (*[]model.User, error) {
	return s.user.FindAllByUsername(username)
}

func (s *userService) FindByUsername(username string) (*model.User, error) {
	return s.user.FindByUsername(username)
}

func (s *userService) FindAllByRoleId(id int64, offset, limit int) (*[]model.User, error) {
	return s.user.FindAllByRoleId(id, offset, limit)
}

func (s *userService) FindUserByEmailOrUsername(user *model.User) (*model.User, error) {
	return s.user.FindUserByEmailOrUsername(user)
}

func (s *userService) UpdateById(uuid *uuid.UUID, user *model.User) (*model.User, error) {
	return s.user.UpdateById(uuid, user)
}

func (s *userService) DeleteById(uuid *uuid.UUID) error {
	return s.user.DeleteById(uuid)
}

func (s *userService) Count() (int64, error) {
	return s.user.Count()
}
