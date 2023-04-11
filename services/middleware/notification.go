package middleware

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	CreatorID    uint
	AssignedToID uint
	Content      string
}

func (n *Notification) Insert() {
	DB.Create(n)
}

func GetAssignedToNotification(assignedToID int) []Notification {
	var notifications []Notification
	DB.Where("assigned_to_id = ?", assignedToID).Find(&notifications)
	return notifications
}
