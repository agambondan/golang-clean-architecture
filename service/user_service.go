package service

import (
	"github.com/google/uuid"
	"golang-youtube-api/models"
	"golang-youtube-api/repository"
)

//var (
//	userRepo repository.UserRepository
//)

type service struct {
	user repository.UserRepository
}

// NewUserService implements the UserService Interface
func NewUserService(repo repository.UserRepository) UserService {
	//userRepo = repo
	return &service{user: repo}
}

//UserApp implements the UserAppInterface
//var _ UserService = &service{}

type UserService interface {
	Validate(user *models.User) error
	Create(user *models.User) (*models.User, error)
	FindAll() ([]models.User, error)
	FindById(uuid uuid.UUID) (models.User, error)
	UpdateById(uuid uuid.UUID, user *models.User) (*models.User, error)
	DeleteById(uuid uuid.UUID) error
}

func (s *service) Validate(user *models.User) error {
	panic("implement me")
}

func (s *service) Create(user *models.User) (*models.User, error) {
	return s.user.Save(user)
}

func (s *service) FindAll() ([]models.User, error) {
	return s.user.FindAll()
	//return userRepo.FindAll()
}

func (s *service) FindById(uuid uuid.UUID) (models.User, error) {
	return s.user.FindById(uuid)
}

func (s *service) UpdateById(uuid uuid.UUID, user *models.User) (*models.User, error) {
	return s.UpdateById(uuid, user)
}

func (s *service) DeleteById(uuid uuid.UUID) error {
	return s.DeleteById(uuid)
}
