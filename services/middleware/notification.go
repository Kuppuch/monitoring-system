package middleware

import "gorm.io/gorm"

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
	DB.Where("assigned_to_id = ?", assignedToID).Find(&notifications)
	return notifications
}

func GetUnreadNotification(assignedToID uint) []Notification {
	var notifications []Notification
	DB.Where("assigned_to_id = ? AND view = false", assignedToID).Find(&notifications)
	return notifications
}
