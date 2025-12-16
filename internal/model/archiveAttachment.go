package model

import "time"

type ArchiveAttachment struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ArchiveHdrID uint       `gorm:"column:archive_hdr_id;not null" json:"archiveHdrId"`
	FileName     string     `gorm:"column:file_name;type:varchar(256);not null" json:"fileName"`
	FileLocation string     `gorm:"column:file_location;type:text;not null" json:"fileLocation"`
	Status       string     `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy    string     `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy   *string    `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt   *time.Time `gorm:"column:modified_at;autoUpdateTime" json:"modifiedAt,omitempty"`

	FileBase64 string `gorm:"-" json:"fileBase64"`
}
