package pkg

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"monitoring-system/config"
)

var DB *gorm.DB

type Mail struct {
	gorm.Model
	Recipient string `json:"recipient"`
	Code      string `json:"code"`
	Content   string `json:"content"`
}

func Connector() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DBPort, config.DbSslMode)
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "notify.", // schema name
			SingularTable: false,
		}})

	//DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&Mail{})
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) InsertMail() (int64, error) {
	tx := DB.Create(&m)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetMailByCode(code, recipient string) Mail {
	m := Mail{}
	DB.Where("code = ? AND recipient = ?", code, recipient).Find(&m)
	return m
}

func DeleteMailByRecipient(recipient string) {
	m := Mail{}
	DB.Where("recipient = ?", recipient).Delete(&m)
}
