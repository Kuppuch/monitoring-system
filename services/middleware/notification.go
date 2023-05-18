package middleware

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"monitoring-system/config"
	"monitoring-system/services/logging"
	"net/http"
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

func (n *Notification) SendSmtpNotify(recipient User) error {
	values := map[string]string{"recipient": recipient.Email, "content": n.Content}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+config.SmtpAddress+"/email/notify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
