package model

import "time"

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       string     `gorm:"column:user_id;type:varchar(256);not null" json:"userId"`
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255)" json:"passwordHash"`
	FullName     string     `gorm:"column:full_name;type:varchar(128)" json:"fullName"`
	DepartmentID uint       `gorm:"column:department_id" json:"departmentId"`
	RoleID       uint       `gorm:"column:role_id" json:"roleId"`
	Status       string     `gorm:"column:status;type:varchar(1);default:'Y'" json:"status"`
	CreatedBy    string     `gorm:"column:created_by;type:varchar(128);not null" json:"createdBy"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ModifiedBy   *string    `gorm:"column:modified_by;type:varchar(128)" json:"modifiedBy,omitempty"`
	ModifiedAt   *time.Time `gorm:"column:modified_at;autoUpdateTime" json:"modifiedAt,omitempty"`

	// Relationships
	Department Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department"`
	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
}
