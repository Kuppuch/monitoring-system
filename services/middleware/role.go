package middleware

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name         string
	ProjectRoles []ProjectRole `gorm:"foreignKey:RoleID" json:"project_roles,omitempty"`
}

type ProjectRole struct {
	gorm.Model
	MemberID uint
	RoleID   uint
}
