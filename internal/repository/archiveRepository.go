package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archive"
	"gorm.io/gorm"
)

type ArchiveRepository interface {
	FindAll() ([]model.ArchiveHdr, error)
	FindByID(id uint) (*model.ArchiveHdr, error)
	Create(archive *model.ArchiveHdr) error
	Update(archive *model.ArchiveHdr) error
	Delete(deleteArchiveRequest *request.DeleteArchiveRequest) error
	FindArchiveByQuery(query string) ([]model.ArchiveHdr, error)
	FindArchiveByAdvanceQuery(advancedSearchRequest request.AdvancedSearchRequest) ([]model.ArchiveHdr, error)
	GetAllArchivesByData(getArchiveByDataRequest request.GetArchiveByDataRequest) ([]model.ArchiveHdr, error)
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
		Preload("ArchiveRoleAccess", "status = ?", "Y").
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
		Preload("ArchiveRoleAccess", "status = ?", "Y").
		Preload("ArchiveRoleAccess.Role").
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

func (r *archiveRepository) Delete(deleteArchiveRequest *request.DeleteArchiveRequest) error {
	result := r.db.Model(&model.ArchiveHdr{}).
		Where("id = ?", deleteArchiveRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteArchiveRequest.SubmittedBy,
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

func (r *archiveRepository) FindArchiveByQuery(queryStr string) ([]model.ArchiveHdr, error) {
	var archives []model.ArchiveHdr

	err := r.db.Model(&model.ArchiveHdr{}).
		Where("status = ?", "Y").
		Where("upper(archive_name) LIKE upper(?)", "%"+queryStr+"%").
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		Preload("ArchiveRoleAccess", "status = ?", "Y").
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
		Preload("ArchiveRoleAccess", "status = ?", "Y").
		Order("archive_name ASC").
		Find(&archives).Error

	return archives, err
}

func (r *archiveRepository) GetAllArchivesByData(getArchiveByDataRequest request.GetArchiveByDataRequest) ([]model.ArchiveHdr, error) {
	var archives []model.ArchiveHdr

	err := r.db.
		Model(&model.ArchiveHdr{}).
		Joins(
			`INNER JOIN archive_role_access
				ON archive_role_access.archive_hdr_id = archive_hdr.id
				AND archive_role_access.role_id = ?
				AND archive_role_access.department_id = ?
			`,
			getArchiveByDataRequest.RoleID,
			getArchiveByDataRequest.DepartmentID,
		).
		Where("archive_role_access.status = ?", "Y").
		Where("archive_hdr.status = ?", "Y").
		Preload("ArchiveAttachments", "status = ?", "Y").
		Preload("ArchiveCharacteristic").
		Preload("ArchiveType").
		Order("archive_date DESC").
		Order("archive_name ASC").
		Find(&archives).Error

	return archives, err
}
