package middleware

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	ProjectID    uint
	UserID       uint
	ProjectRoles []ProjectRole `gorm:"foreignKey:MemberID" json:"project_roles,omitempty"`
}
