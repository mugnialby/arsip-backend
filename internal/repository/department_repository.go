package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/department"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	FindAll() ([]model.Department, error)
	FindByID(id uint) (*model.Department, error)
	Create(department *model.Department) error
	Update(department *model.Department) error
	Delete(deleteDepartmentRequest *request.DeleteDepartmentRequest) error
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

func (r *departmentRepository) Delete(deleteDepartmentRequest *request.DeleteDepartmentRequest) error {
	result := r.db.Model(&model.Department{}).
		Where("id = ?", deleteDepartmentRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteDepartmentRequest.SubmittedBy,
			"modified_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no data found to delete")
	}

	return nil
}
