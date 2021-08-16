package service

import (
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
)

type roleService struct {
	role repository.RoleRepository
}

// NewRoleService implements the RoleService Interface
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo}
}

type RoleService interface {
	Create(role *model.Role) (*model.Role, error)
	FindAll() ([]model.Role, error)
	FindById(id uint64) (model.Role, error)
	UpdateById(id uint64, role *model.Role) (*model.Role, error)
	DeleteById(id uint64) error
	Count() (int, error)
}

func (p *roleService) Create(role *model.Role) (*model.Role, error) {
	return p.role.Save(role)
}

func (p *roleService) FindAll() ([]model.Role, error) {
	return p.role.FindAll()
}

func (p *roleService) FindById(id uint64) (model.Role, error) {
	return p.role.FindById(id)
}

func (p *roleService) UpdateById(id uint64, role *model.Role) (*model.Role, error) {
	return p.role.UpdateById(id, role)
}

func (p *roleService) DeleteById(id uint64) error {
	return p.role.DeleteById(id)
}

func (p *roleService) Count() (int,error) {
	return p.role.Count()
}