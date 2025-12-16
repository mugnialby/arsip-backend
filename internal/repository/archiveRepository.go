package repository

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archive"
	"gorm.io/gorm"
)

type ArchiveRepository interface {
	FindAll() ([]model.ArchiveHdr, error)
	FindByID(id uint) (*model.ArchiveHdr, error)
	Create(archive *model.ArchiveHdr) error
	Update(archive *model.ArchiveHdr) error
	FindArchiveByQuery(query string) ([]model.ArchiveHdr, error)
	FindArchiveByAdvanceQuery(advancedSearchRequest request.AdvancedSearchRequest) ([]model.ArchiveHdr, error)
}

type archiveRepository struct {
	db *gorm.DB
}

func NewArchiveRepository(db *gorm.DB) ArchiveRepository {
	return &archiveRepository{db: db}
}

func (r *archiveRepository) FindAll() ([]model.ArchiveHdr, error) {
	var archives []model.ArchiveHdr

	err := r.db.Model(&model.ArchiveHdr{}).
		Where("status = ?", "Y").
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		Order("archive_date ASC").
		Find(&archives).Error

	return archives, err
}

func (r *archiveRepository) FindByID(id uint) (*model.ArchiveHdr, error) {
	var archive model.ArchiveHdr

	err := r.db.Model(&model.ArchiveHdr{}).
		Where("id = ? AND status = ?", id, "Y").
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		First(&archive).Error

	return &archive, err
}

func (r *archiveRepository) Create(archive *model.ArchiveHdr) error {
	return r.db.Create(archive).Error
}

func (r *archiveRepository) Update(archive *model.ArchiveHdr) error {
	return r.db.Model(&model.ArchiveHdr{}).
		Where("id = ?", archive.ID).
		Updates(archive).Error
}

func (r *archiveRepository) FindArchiveByQuery(queryStr string) ([]model.ArchiveHdr, error) {
	var archives []model.ArchiveHdr

	err := r.db.Model(&model.ArchiveHdr{}).
		Where("status = ?", "Y").
		Where("upper(archive_name) LIKE upper(?)", "%"+queryStr+"%").
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		Order("archive_name ASC").
		Find(&archives).Error

	return archives, err
}

func (r *archiveRepository) FindArchiveByAdvanceQuery(req request.AdvancedSearchRequest) ([]model.ArchiveHdr, error) {
	var archives []model.ArchiveHdr

	q := r.db.Model(&model.ArchiveHdr{}).
		Where("status = ?", "Y")

	if req.ArchiveName != nil && *req.ArchiveName != "" {
		q = q.Where("upper(archive_name) LIKE upper(?)", "%"+*req.ArchiveName+"%")
	}

	err := q.
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		Order("archive_name ASC").
		Find(&archives).Error

	return archives, err
}
