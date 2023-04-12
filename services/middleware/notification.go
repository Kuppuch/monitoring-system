package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type Notification struct {
	gorm.Model
	View         bool
	CreatorID    uint
	AssignedToID uint
	Content      string
	Source       string
}

func (n *Notification) Insert() {
	DB.Create(n)
}

func GetAssignedToNotification(assignedToID uint) []Notification {
	var notifications []Notification
	DB.Where("assigned_to_id = ?", assignedToID).Order("id desc").Find(&notifications)
	return notifications
}

func GetUnreadNotification(assignedToID uint) []Notification {
	var notifications []Notification
	DB.Where("assigned_to_id = ? AND view = false", assignedToID).Find(&notifications)
	return notifications
}

func SetReadNotification(id int) (int, error) {
	tx := DB.Table("notifications").Where("id = ?", id).Update("view", "true")
	if tx.Error != nil {
		logging.Print.Error(tx.Error)
		return 0, tx.Error
	}
	return int(tx.RowsAffected), nil
}
