package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Member struct {
	gorm.Model
	ProjectID    uint
	UserID       uint
	ProjectRoles []ProjectRole `gorm:"foreignKey:MemberID" json:"project_roles,omitempty"`
}

type MemberView struct {
	gorm.Model
	Name     string
	LastName string
	Role     string
}

func (m *Member) GetMember() {
	DB.Where("project_id = ? AND user_id = ?", m.ProjectID, m.UserID).Find(m)
}

func (m *Member) InsertMember() int64 {
	tx := DB.Create(m)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func GetMembers(projectId int) []MemberView {
	memberView := []MemberView{}
	tx := DB.Table("project_roles as pr").Select("u.name as name, u.last_name as last_name, r.name as role").
		Joins("INNER JOIN members as m ON m.id = pr.member_id").
		Joins("INNER JOIN roles as r ON r.id = pr.role_id").
		Joins("INNER JOIN users as u ON u.id = m.user_id").
		Where("m.project_id = ?", projectId).Find(&memberView)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return memberView
}
