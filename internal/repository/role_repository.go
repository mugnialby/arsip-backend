package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/roles"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindAll() ([]model.Role, error)
	FindByID(id uint) (*model.Role, error)
	Create(book *model.Role) error
	Update(book *model.Role) error
	Delete(deleteRoleRequest *request.DeleteRoleRequest) error
	GetRoleByDepartmentID(departmentId uint) ([]model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role

	err := r.db.
		Model(&model.Role{}).
		Joins("INNER JOIN departments ON departments.id = roles.department_id").
		Where("roles.status = ?", "Y").
		Where("roles.id <> ?", 1).
		Preload("Department", "status = ?", "Y").
		Order("departments.department_name ASC").
		Order("roles.role_name ASC").
		Find(&roles).Error

	return roles, err

}

func (r *roleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(deleteRoleRequest *request.DeleteRoleRequest) error {
	result := r.db.Model(&model.Role{}).
		Where("id = ?", deleteRoleRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteRoleRequest.SubmittedBy,
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

func (r *roleRepository) GetRoleByDepartmentID(departmentId uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Where("status = ?", "Y").
		Where("department_id = ?", departmentId).
		Order("department_id asc").
		Order("role_name asc").
		Find(&roles).Error
	return roles, err
}
