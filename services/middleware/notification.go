package middleware

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"monitoring-system/config"
	"monitoring-system/services/api/socket"
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

func (n *Notification) Prepare(issue IssueWeb) {

	n.Insert()
	// Отправляем уведомление в канал уведомлений
	nByte, _ := json.Marshal(&n)
	socket.BigChannel <- nByte
	assignedUser := User{
		Model: gorm.Model{ID: n.AssignedToID},
	}
	if err := assignedUser.GetUser(); err != nil {
		logging.Print.Error("failed get assigned user for notify issue create by email ", err)
		return
	}
	if assignedUser.EmailNotify && len(config.SmtpAddress) > 0 {
		if err := n.SendSmtpNotify(assignedUser); err != nil {
			logging.Print.Error("failed send email to assigned for create issue ", err)
			return
		}
	}
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
