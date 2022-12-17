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
