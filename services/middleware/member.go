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
