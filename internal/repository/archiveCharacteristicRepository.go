package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"gorm.io/gorm"
)

type ArchiveCharacteristicRepository interface {
	FindAll() ([]model.ArchiveCharacteristic, error)
	FindByID(id uint) (*model.ArchiveCharacteristic, error)
	Create(archiveCharacteristic *model.ArchiveCharacteristic) error
	Update(archiveCharacteristic *model.ArchiveCharacteristic) error
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
