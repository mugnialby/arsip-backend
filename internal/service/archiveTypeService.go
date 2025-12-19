package service

import (
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveType"
	"github.com/mugnialby/arsip-backend/internal/repository"
)

type ArchiveTypeService struct {
	repo repository.ArchiveTypeRepository
}

func NewArchiveTypeService(repo repository.ArchiveTypeRepository) *ArchiveTypeService {
	return &ArchiveTypeService{repo: repo}
}

func (s *ArchiveTypeService) GetAllArchiveTypes() ([]model.ArchiveType, error) {
	return s.repo.FindAll()
}

func (s *ArchiveTypeService) GetArchiveTypeByID(id uint) (*model.ArchiveType, error) {
	return s.repo.FindByID(id)
}

func (s *ArchiveTypeService) CreateArchiveType(archiveType *model.ArchiveType) error {
	return s.repo.Create(archiveType)
}

func (s *ArchiveTypeService) UpdateArchiveType(archiveType *model.ArchiveType) error {
	return s.repo.Update(archiveType)
}

func (s *ArchiveTypeService) DeleteArchiveType(deleteArchiveTypeRequest *request.DeleteArchiveTypeRequest) error {
	return s.repo.Delete(deleteArchiveTypeRequest)
}
