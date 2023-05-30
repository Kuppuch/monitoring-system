package middleware

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"monitoring-system/services/api/socket"
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
	tx := DB.Where("assigned_to_id = ? AND view = false", assignedToID).Find(&notifications)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
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
	if assignedUser.EmailNotify {
		if err := n.SendSmtpNotify(assignedUser); err != nil {
			logging.Print.Error("failed send email to assigned for create issue ", err)
			return
		}
	}
}

func (n *Notification) SendSmtpNotify(recipient User) error {

	err := sender(recipient.Email, "", n.Content)
	if err != nil {
		return err
	}
	return nil
}

func sender(recipient string, code string, content string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "vehicleaggregator@gmail.com")
	mail.SetHeader("To", recipient)

	if len(code) > 0 {
		mail.SetHeader("Subject", "Confirm action")
		mail.SetBody("text/html", fmt.Sprintf("Your verification code: %s", code))
	} else {
		mail.SetHeader("Subject", "Уведомление системы мониторинга")
		mail.SetBody("text/html", content)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, "vehicleaggregator@gmail.com", "hcbwxhrlwxgovdtu")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
