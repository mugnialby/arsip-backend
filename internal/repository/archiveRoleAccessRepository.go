package repository

import (
	"errors"
	"time"

	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveRoleAccess"
	"gorm.io/gorm"
)

type ArchiveRoleAccessRepository interface {
	FindAll() ([]model.ArchiveRoleAccess, error)
	FindByID(id uint) (*model.ArchiveRoleAccess, error)
	Create(archiveRoleAccess *model.ArchiveRoleAccess) error
	Update(archiveRoleAccess *model.ArchiveRoleAccess) error
	Delete(deleteArchiveRoleAccessRequest *request.DeleteArchiveRoleAccessRequest) error
	DeleteArchiveRoleAccessByArchiveID(archiveID uint, submittedBy string) error
}

type archiveRoleAccessRepository struct {
	db *gorm.DB
}

func NewArchiveRoleAccessRepository(db *gorm.DB) ArchiveRoleAccessRepository {
	return &archiveRoleAccessRepository{db: db}
}

func (r *archiveRoleAccessRepository) FindAll() ([]model.ArchiveRoleAccess, error) {
	var books []model.ArchiveRoleAccess
	err := r.db.Where("status = ?", "Y").
		Find(&books).Error
	return books, err
}

func (r *archiveRoleAccessRepository) FindByID(id uint) (*model.ArchiveRoleAccess, error) {
	var archiveRoleAccess model.ArchiveRoleAccess
	err := r.db.First(&archiveRoleAccess, id).Error
	return &archiveRoleAccess, err
}

func (r *archiveRoleAccessRepository) Create(archiveRoleAccess *model.ArchiveRoleAccess) error {
	return r.db.Create(archiveRoleAccess).Error
}

func (r *archiveRoleAccessRepository) Update(archiveRoleAccess *model.ArchiveRoleAccess) error {
	return r.db.Save(archiveRoleAccess).Error
}

func (r *archiveRoleAccessRepository) Delete(deleteArchiveRoleAccessRequest *request.DeleteArchiveRoleAccessRequest) error {
	result := r.db.Model(&model.ArchiveRoleAccess{}).
		Where("id = ?", deleteArchiveRoleAccessRequest.ID).
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": deleteArchiveRoleAccessRequest.SubmittedBy,
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

func (r *archiveRoleAccessRepository) DeleteArchiveRoleAccessByArchiveID(archiveID uint, submittedBy string) error {
	result := r.db.Model(&model.ArchiveRoleAccess{}).
		Where("archive_hdr_id = ?", archiveID).
		Where("status = ?", "Y").
		Updates(map[string]interface{}{
			"status":      "N",
			"modified_by": submittedBy,
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
