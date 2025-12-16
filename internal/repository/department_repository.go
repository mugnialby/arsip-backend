package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	FindAll() ([]model.Department, error)
	FindByID(id uint) (*model.Department, error)
	Create(department *model.Department) error
	Update(department *model.Department) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) FindAll() ([]model.Department, error) {
	var departments []model.Department
	err := r.db.Where("status = ?", "Y").
		Order("department_name asc").
		Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) FindByID(id uint) (*model.Department, error) {
	var department model.Department
	err := r.db.First(&department, id).Error
	return &department, err
}

func (r *departmentRepository) Create(department *model.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) Update(department *model.Department) error {
	return r.db.Save(department).Error
}
