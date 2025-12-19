package model

import (
	"time"
)

type ArchiveRoleAccess struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ArchiveHdrID uint       `gorm:"column:archive_hdr_id" json:"archiveHdrId"`
	RoleID       uint       `gorm:"column:role_id" json:"roleId"`
	DepartmentID uint       `gorm:"column:department_id" json:"departmentId"`
	Status       string     `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy    string     `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy   *string    `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt   *time.Time `gorm:"column:modified_at;" json:"modifiedAt,omitempty"`

	Role *Role `gorm:"foreignKey:RoleID;->" json:"role"`
}

func (ArchiveRoleAccess) TableName() string {
	return "archive_role_access"
}
