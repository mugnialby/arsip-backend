package model

import (
	"time"

	"github.com/mugnialby/arsip-backend/internal/utils"
)

type ArchiveHdr struct {
	ID                      uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ArchiveName             string         `gorm:"column:archive_name;type:varchar(256);not null" json:"archiveName"`
	ArchiveNumber           string         `gorm:"column:archive_number;type:varchar(256);not null" json:"archiveNumber"`
	ArchiveCharacteristicID uint           `gorm:"column:archive_characteristic_id" json:"archiveCharacteristicId"`
	ArchiveTypeID           uint           `gorm:"column:archive_type_id" json:"archiveTypeId"`
	ArchiveDate             utils.DateOnly `gorm:"type:date column:archive_date" json:"archiveDate"`
	DepartmentID            uint           `gorm:"column:department_id" json:"departmentId"`
	Status                  string         `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy               string         `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt               time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy              *string        `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt              *time.Time     `gorm:"column:modified_at;" json:"modifiedAt,omitempty"`

	// Read-only relations
	ArchiveCharacteristic *ArchiveCharacteristic `gorm:"foreignKey:ArchiveCharacteristicID;->" json:"archiveCharacteristic"`
	ArchiveType           *ArchiveType           `gorm:"foreignKey:ArchiveTypeID;->" json:"archiveType"`
	Department            *Department            `gorm:"foreignKey:DepartmentID;->" json:"department"`

	ArchiveRoleAccess  []*ArchiveRoleAccess `gorm:"foreignKey:ArchiveHdrID;->" json:"archiveRoleAccess"`
	ArchiveAttachments []*ArchiveAttachment `gorm:"foreignKey:ArchiveHdrID;->" json:"archiveAttachments"`
}

func (ArchiveHdr) TableName() string {
	return "archive_hdr"
}
