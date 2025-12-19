package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveType"
	"gorm.io/gorm"
)

type ArchiveTypeRepository interface {
	FindAll() ([]model.ArchiveType, error)
	FindByID(id uint) (*model.ArchiveType, error)
	Create(archiveType *model.ArchiveType) error
	Delete(deleteArchiveTypeRequest *request.DeleteArchiveTypeRequest) error
	Update(archiveType *model.ArchiveType) error
}

type archiveTypeRepository struct {
	db *gorm.DB
}

func NewArchiveTypeRepository(db *gorm.DB) ArchiveTypeRepository {
	return &archiveTypeRepository{db: db}
}

func (r *archiveTypeRepository) FindAll() ([]model.ArchiveType, error) {
	var archiveTypes []model.ArchiveType
	err := r.db.Where("status = ?", "Y").
		Order("archive_type_name asc").
		Find(&archiveTypes).Error
	return archiveTypes, err
}

func (r *archiveTypeRepository) FindByID(id uint) (*model.ArchiveType, error) {
	var archiveType model.ArchiveType
	err := r.db.First(&archiveType, id).Error
	return &archiveType, err
}

func (r *archiveTypeRepository) Create(archiveType *model.ArchiveType) error {
	return r.db.Create(archiveType).Error
}

func (r *archiveTypeRepository) Update(archiveType *model.ArchiveType) error {
	return r.db.Save(archiveType).Error
}

func (r *archiveTypeRepository) Delete(deleteArchiveTypeRequest *request.DeleteArchiveTypeRequest) error {
	result := r.db.Model(&model.ArchiveType{}).
		Where("id = ?", deleteArchiveTypeRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteArchiveTypeRequest.SubmittedBy,
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
