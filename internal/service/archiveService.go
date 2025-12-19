package service

import (
	"github.com/mugnialby/arsip-backend/internal/model"
	request "github.com/mugnialby/arsip-backend/internal/model/dto/request/archive"
	"github.com/mugnialby/arsip-backend/internal/repository"
)

type ArchiveService struct {
	repo repository.ArchiveRepository
}

func NewArchiveService(repo repository.ArchiveRepository) *ArchiveService {
	return &ArchiveService{repo: repo}
}

func (s *ArchiveService) GetAllArchives() ([]model.ArchiveHdr, error) {
	return s.repo.FindAll()
}

func (s *ArchiveService) GetArchiveByID(id uint) (*model.ArchiveHdr, error) {
	return s.repo.FindByID(id)
}

func (s *ArchiveService) CreateArchive(archive *model.ArchiveHdr) error {
	return s.repo.Create(archive)
}

func (s *ArchiveService) UpdateArchive(archive *model.ArchiveHdr) error {
	return s.repo.Update(archive)
}

func (s *ArchiveService) FindArchiveByQuery(query string) ([]model.ArchiveHdr, error) {
	return s.repo.FindArchiveByQuery(query)
}

func (s *ArchiveService) DeleteArchive(deleteArchiveRequest *request.DeleteArchiveRequest) error {
	return s.repo.Delete(deleteArchiveRequest)
}

func (s *ArchiveService) FindArchiveByAdvanceQuery(advancedSearchRequest request.AdvancedSearchRequest) ([]model.ArchiveHdr, error) {
	return s.repo.FindArchiveByAdvanceQuery(advancedSearchRequest)
}

func (s *ArchiveService) GetAllArchivesByData(getArchiveByDataRequest request.GetArchiveByDataRequest) ([]model.ArchiveHdr, error) {
	return s.repo.GetAllArchivesByData(getArchiveByDataRequest)
}
