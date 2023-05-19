package middleware

import "gorm.io/gorm"

type Observe struct {
	gorm.Model
	IssueID int
	UserID  int
}

func (o *Observe) Insert() error {
	tx := DB.Create(o)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetSubscribers(issueID uint) []User {
	var users []User
	DB.Model(&User{}).Joins("INNER JOIN observes AS o ON o.user_id = users.id").
		Where("o.issue_id = ?", issueID).Find(&users)
	return users
}
