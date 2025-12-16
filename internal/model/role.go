package model

import "time"

type Role struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentID uint       `gorm:"column:department_id;type:varchar(128);not null" json:"departmentId"`
	RoleName     string     `gorm:"column:role_name;type:varchar(128);not null" json:"roleName"`
	Status       string     `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy    string     `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy   *string    `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt   *time.Time `gorm:"column:modified_at;autoUpdateTime" json:"modifiedAt,omitempty"`

	// Relationships
	Department Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department"`
}
