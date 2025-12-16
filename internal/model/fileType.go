package model

import "time"

type FileType struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FileTypeName string     `gorm:"column:file_type_name;type:varchar(128);not null" json:"fileTypeName"`
	Status       string     `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy    string     `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy   *string    `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt   *time.Time `gorm:"column:modified_at;autoUpdateTime" json:"modifiedAt,omitempty"`
}
