package service

import (
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
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
	FindAll() (*[]model.Role, error)
	FindById(id *int64) (*model.Role, error)
	UpdateById(id *int64, role *model.Role) (*model.Role, error)
	DeleteById(id *int64) error
	Count() (*int64, error)
}

func (p *roleService) Create(role *model.Role) (*model.Role, error) {
	return p.role.Save(role)
}

func (p *roleService) FindAll() (*[]model.Role, error) {
	return p.role.FindAll()
}

func (p *roleService) FindById(id *int64) (*model.Role, error) {
	return p.role.FindById(id)
}

func (p *roleService) UpdateById(id *int64, role *model.Role) (*model.Role, error) {
	return p.role.UpdateById(id, role)
}

func (p *roleService) DeleteById(id *int64) error {
	return p.role.DeleteById(id)
}

func (p *roleService) Count() (*int64, error) {
	return p.role.Count()
}
