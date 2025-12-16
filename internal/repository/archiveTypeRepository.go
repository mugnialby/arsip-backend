package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"gorm.io/gorm"
)

type ArchiveTypeRepository interface {
	FindAll() ([]model.ArchiveType, error)
	FindByID(id uint) (*model.ArchiveType, error)
	Create(archiveType *model.ArchiveType) error
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
