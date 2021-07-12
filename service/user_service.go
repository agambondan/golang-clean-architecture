package service

import (
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
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
	FindAll() ([]model.User, error)
	FindById(uuid uuid.UUID) (model.User, error)
	FindAllByRoleId(id uint64) ([]model.User, error)
	FindUserByEmailAndPassword(user *model.User) (model.User, error)
	UpdateById(uuid uuid.UUID, user *model.User) (*model.User, error)
	DeleteById(uuid uuid.UUID) error
}

func (s *userService) Create(user *model.User) (*model.User, error) {
	return s.user.Save(user)
}

func (s *userService) FindAll() ([]model.User, error) {
	return s.user.FindAll()
}

func (s *userService) FindById(uuid uuid.UUID) (model.User, error) {
	return s.user.FindById(uuid)
}

func (s *userService) FindAllByRoleId(id uint64) ([]model.User, error) {
	return s.user.FindAllByRoleId(id)
}

func (s *userService) FindUserByEmailAndPassword(user *model.User) (model.User, error) {
	return s.user.FindUserByEmailAndPassword(user)
}

func (s *userService) UpdateById(uuid uuid.UUID, user *model.User) (*model.User, error) {
	return s.user.UpdateById(uuid, user)
}

func (s *userService) DeleteById(uuid uuid.UUID) error {
	return s.user.DeleteById(uuid)
}
