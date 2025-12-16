package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type ArchiveAttachmentService struct {
	repo repository.ArchiveAttachmentRepository
}

func NewArchiveAttachmentService(repo repository.ArchiveAttachmentRepository) *ArchiveAttachmentService {
	return &ArchiveAttachmentService{repo: repo}
}

func (s *ArchiveAttachmentService) GetAllArchiveAttachments() ([]model.ArchiveAttachment, error) {
	return s.repo.FindAll()
}

func (s *ArchiveAttachmentService) GetArchiveAttachmentByID(id uint) (*model.ArchiveAttachment, error) {
	return s.repo.FindByID(id)
}

func (s *ArchiveAttachmentService) CreateArchiveAttachment(archiveAttachment *model.ArchiveAttachment) error {
	return s.repo.Create(archiveAttachment)
}

func (s *ArchiveAttachmentService) UpdateArchiveAttachment(archiveAttachment *model.ArchiveAttachment) error {
	return s.repo.Update(archiveAttachment)
}

func (s *ArchiveAttachmentService) DeleteArchiveAttachmentByArchiveID(archiveID uint) error {
	return s.repo.DeleteArchiveAttachmentByArchiveID(archiveID)
}
