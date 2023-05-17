package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"monitoring-system/notify_service/pkg"
	"net/http"
	"strconv"
	"time"
)

func SendEmail(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed get raw data",
		})
		return
	}

	m := pkg.Mail{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed unmarshal raw data",
		})
		return
	}
	m.Code = strconv.Itoa(rand.Intn(8999) + 1000)
	err = sender(m)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed send email",
		})
		return
	}
	row, err := m.InsertMail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed insert code into database",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"row_affected": row})
}

func GetCodeByEmail(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed get raw data",
		})
		return
	}

	m := pkg.Mail{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed unmarshal raw data",
		})
		return
	}

	mm := pkg.GetMailByCode(m.Code, m.Recipient)
	after := mm.CreatedAt.Add(10 * time.Minute).After(time.Now())
	if !after || mm.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"valid": false})
		return
	}
	pkg.DeleteMailByRecipient(mm.Recipient)
	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func SendEmailNotify(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed get raw data",
		})
		return
	}

	m := pkg.Mail{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed unmarshal raw data",
		})
		return
	}
	err = sender(m)
}

func sender(m pkg.Mail) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "vehicleaggregator@gmail.com")
	mail.SetHeader("To", m.Recipient)

	//rand.Seed(time.Now().UnixNano())
	if len(m.Code) > 0 {
		mail.SetHeader("Subject", "Confirm action")
		mail.SetBody("text/html", fmt.Sprintf("Your verification code: %s", m.Code))
	} else {
		mail.SetHeader("Subject", "Уведомление системы мониторинга")
		mail.SetBody("text/html", m.Content)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, "vehicleaggregator@gmail.com", "hcbwxhrlwxgovdtu")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
