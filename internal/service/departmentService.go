package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/department"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type DepartmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) GetAllDepartments() ([]model.Department, error) {
	return s.repo.FindAll()
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*model.Department, error) {
	return s.repo.FindByID(id)
}

func (s *DepartmentService) CreateDepartment(department *model.Department) error {
	return s.repo.Create(department)
}

func (s *DepartmentService) UpdateDepartment(department *model.Department) error {
	return s.repo.Update(department)
}

func (s *DepartmentService) DeleteDepartment(deleteDepartmentRequest *request.DeleteDepartmentRequest) error {
	return s.repo.Delete(deleteDepartmentRequest)
}
