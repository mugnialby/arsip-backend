package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveRoleAccess"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type ArchiveRoleAccessService struct {
	repo repository.ArchiveRoleAccessRepository
}

func NewArchiveRoleAccessService(repo repository.ArchiveRoleAccessRepository) *ArchiveRoleAccessService {
	return &ArchiveRoleAccessService{repo: repo}
}

func (s *ArchiveRoleAccessService) GetAllArchiveRoleAccesss() ([]model.ArchiveRoleAccess, error) {
	return s.repo.FindAll()
}

func (s *ArchiveRoleAccessService) GetArchiveRoleAccessByID(id uint) (*model.ArchiveRoleAccess, error) {
	return s.repo.FindByID(id)
}

func (s *ArchiveRoleAccessService) CreateArchiveRoleAccess(archiveRoleAccess *model.ArchiveRoleAccess) error {
	return s.repo.Create(archiveRoleAccess)
}

func (s *ArchiveRoleAccessService) UpdateArchiveRoleAccess(archiveRoleAccess *model.ArchiveRoleAccess) error {
	return s.repo.Update(archiveRoleAccess)
}

func (s *ArchiveRoleAccessService) DeleteArchiveRoleAccess(deleteRoleAccessRequest *request.DeleteArchiveRoleAccessRequest) error {
	return s.repo.Delete(deleteRoleAccessRequest)
}

func (s *ArchiveRoleAccessService) DeleteArchiveRoleAccessByArchiveID(archiveID uint, submittedBy string) error {
	return s.repo.DeleteArchiveRoleAccessByArchiveID(archiveID, submittedBy)
}
