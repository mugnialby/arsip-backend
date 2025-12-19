package service

import (
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model"
	request "github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/model/dto/request/archiveCharacteristic"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/repository"
)

type ArchiveCharacteristicService struct {
	repo repository.ArchiveCharacteristicRepository
}

func NewArchiveCharacteristicService(repo repository.ArchiveCharacteristicRepository) *ArchiveCharacteristicService {
	return &ArchiveCharacteristicService{repo: repo}
}

func (s *ArchiveCharacteristicService) GetAllArchiveCharacteristics() ([]model.ArchiveCharacteristic, error) {
	return s.repo.FindAll()
}

func (s *ArchiveCharacteristicService) GetArchiveCharacteristicByID(id uint) (*model.ArchiveCharacteristic, error) {
	return s.repo.FindByID(id)
}

func (s *ArchiveCharacteristicService) CreateArchiveCharacteristic(archiveCharacteristic *model.ArchiveCharacteristic) error {
	return s.repo.Create(archiveCharacteristic)
}

func (s *ArchiveCharacteristicService) UpdateArchiveCharacteristic(archiveCharacteristic *model.ArchiveCharacteristic) error {
	return s.repo.Update(archiveCharacteristic)
}

func (s *ArchiveCharacteristicService) DeleteArchiveCharacteristic(deleteArchiveCharacteristicRequest *request.DeleteArchiveCharacteristicRequest) error {
	return s.repo.Delete(deleteArchiveCharacteristicRequest)
}
