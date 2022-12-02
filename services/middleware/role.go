package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

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

func GetRoles() []Role {
	var roles []Role
	DB.Find(&roles)
	return roles
}

func (r *ProjectRole) GetProjectRole() {
	DB.Where("member_id = ? AND role_id = ?", r.MemberID, r.RoleID).Find(r)
}

func (r *ProjectRole) InsertRole() int64 {
	tx := DB.Create(r)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}
