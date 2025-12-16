package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type RoleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	return s.repo.FindAll()
}

func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	return s.repo.FindByID(id)
}

func (s *RoleService) CreateRole(role *model.Role) error {
	return s.repo.Create(role)
}

func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.repo.Update(role)
}

func (s *RoleService) GetRoleByDepartmentID(departmentId uint) ([]model.Role, error) {
	return s.repo.GetRoleByDepartmentID(departmentId)
}
