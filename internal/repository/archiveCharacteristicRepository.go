package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveCharacteristic"
	"gorm.io/gorm"
)

type ArchiveCharacteristicRepository interface {
	FindAll() ([]model.ArchiveCharacteristic, error)
	FindByID(id uint) (*model.ArchiveCharacteristic, error)
	Create(archiveCharacteristic *model.ArchiveCharacteristic) error
	Update(archiveCharacteristic *model.ArchiveCharacteristic) error
	Delete(deleteUserRequest *request.DeleteArchiveCharacteristicRequest) error
}

type archiveCharacteristicRepository struct {
	db *gorm.DB
}

func NewArchiveCharacteristicRepository(db *gorm.DB) ArchiveCharacteristicRepository {
	return &archiveCharacteristicRepository{db: db}
}

func (r *archiveCharacteristicRepository) FindAll() ([]model.ArchiveCharacteristic, error) {
	var archiveCharacteristics []model.ArchiveCharacteristic
	err := r.db.Where("status = ?", "Y").
		Order("archive_characteristic_name asc").
		Find(&archiveCharacteristics).Error
	return archiveCharacteristics, err
}

func (r *archiveCharacteristicRepository) FindByID(id uint) (*model.ArchiveCharacteristic, error) {
	var archiveCharacteristic model.ArchiveCharacteristic
	err := r.db.First(&archiveCharacteristic, id).Error
	return &archiveCharacteristic, err
}

func (r *archiveCharacteristicRepository) Create(archiveCharacteristic *model.ArchiveCharacteristic) error {
	return r.db.Create(archiveCharacteristic).Error
}

func (r *archiveCharacteristicRepository) Update(archiveCharacteristic *model.ArchiveCharacteristic) error {
	return r.db.Save(archiveCharacteristic).Error
}

func (r *archiveCharacteristicRepository) Delete(deleteArchiveCharacteristicRequest *request.DeleteArchiveCharacteristicRequest) error {
	result := r.db.Model(&model.ArchiveCharacteristic{}).
		Where("id = ?", deleteArchiveCharacteristicRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteArchiveCharacteristicRequest.SubmittedBy,
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
