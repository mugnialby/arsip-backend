package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"gorm.io/gorm"
)

type ArchiveAttachmentRepository interface {
	FindAll() ([]model.ArchiveAttachment, error)
	FindByID(id uint) (*model.ArchiveAttachment, error)
	Create(archiveAttachment *model.ArchiveAttachment) error
	Update(archiveAttachment *model.ArchiveAttachment) error
	DeleteArchiveAttachmentByArchiveID(archiveID uint) error
}

type archiveAttachmentRepository struct {
	db *gorm.DB
}

func NewArchiveAttachmentRepository(db *gorm.DB) ArchiveAttachmentRepository {
	return &archiveAttachmentRepository{db: db}
}

func (r *archiveAttachmentRepository) FindAll() ([]model.ArchiveAttachment, error) {
	var books []model.ArchiveAttachment
	err := r.db.Where("status = ?", "Y").
		Find(&books).Error
	return books, err
}

func (r *archiveAttachmentRepository) FindByID(id uint) (*model.ArchiveAttachment, error) {
	var archiveAttachment model.ArchiveAttachment
	err := r.db.First(&archiveAttachment, id).Error
	return &archiveAttachment, err
}

func (r *archiveAttachmentRepository) Create(archiveAttachment *model.ArchiveAttachment) error {
	return r.db.Create(archiveAttachment).Error
}

func (r *archiveAttachmentRepository) Update(archiveAttachment *model.ArchiveAttachment) error {
	return r.db.Save(archiveAttachment).Error
}

func (r *archiveAttachmentRepository) DeleteArchiveAttachmentByArchiveID(archiveID uint) error {
	return r.db.Model(&model.ArchiveAttachment{}).
		Where("archive_hdr_id = ? AND status = ?", archiveID, "Y").
		Update("status", "N").Error
}
