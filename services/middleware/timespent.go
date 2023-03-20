package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Timespent struct {
	gorm.Model
	IssueID uint    `json:"issue_id"`
	UserID  uint    `json:"user_id"`
	RoleID  uint    `json:"role_id"`
	Spent   float32 `json:"spent"`
}

func (timespent *Timespent) Insert() (int64, error) {
	tx := DB.Create(&timespent)
	if tx.Error != nil {
		logging.Print.Error("database error create timespent")
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetUserIssueTimespent(issueID int, userID int) []Timespent {
	var t []Timespent
	DB.Where("issue_id = ? AND user_id = ?", issueID, userID).Find(&t)
	return t
}

type TimespentReport struct {
	ID        int     `json:"ID"`
	Spent     float32 `json:"spent"`
	Issue     string  `json:"issue_name"`
	Developer string  `json:"developer"`
	Role      string  `json:"role"`
	RoleID    int     `json:"role_id"`
}

func GetBudgetTimespent(id int) []TimespentReport {
	var timespents []TimespentReport
	tx := DB.Table("timespents as t").
		Select("t.id, i.name AS issue, u.name || ' ' || u.last_name AS Developer, t.spent, r.name AS Role, r.id AS Role_id").
		Joins("INNER JOIN issues AS i ON t.issue_id = i.id").
		Joins("INNER JOIN users AS u ON u.id = t.user_id").
		Joins("LEFT JOIN roles AS r ON r.id = t.role_id").
		Where("i.budget_id = ?", id).
		Order("t.id").Find(&timespents)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return timespents
}
