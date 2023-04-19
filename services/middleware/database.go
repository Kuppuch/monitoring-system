package middleware

import (
	"crypto/sha256"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"monitoring-system/config"
	"monitoring-system/services/logging"
	"os"
	"strconv"
)

var DB *gorm.DB

func Connector() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DBPort, config.DbSslMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&User{}, &Status{}, &Tracker{}, &Budget{}, &Issue{}, &Project{}, &Auth{}, &Member{}, &Role{},
		&ProjectRole{}, &Timespent{}, &GitRepository{}, &Notification{}, &Risk{}, &HeadRole{})
	if err != nil {
		return err
	}
	insertStatuses()
	insertTrackers()
	insertAdmin()
	insertRoles()

	return nil
}

func insertStatuses() {
	var cas []Status
	DB.Find(&cas)
	if len(cas) == 0 {
		data := []string{
			"Новая",
			"В очереди исполнения",
			"В работе",
			"В тестировании",
			"Возвращено в работу",
			"Готова",
			"Переоткрыта",
			"Закрыта",
		}
		for _, v := range data {
			DB.Create(&Status{
				Name: v,
			})
		}
	}
}

func insertTrackers() {
	var trackers []Tracker
	DB.Find(&trackers)
	if len(trackers) == 0 {
		data := []string{
			"Фича",
			"Анализ",
			"Баг",
		}
		for _, v := range data {
			DB.Create(&Tracker{
				Name: v,
			})
		}
	}
}

func insertAdmin() {
	var user []User
	DB.Where("email = ?", "admin@admin.ru").Find(&user)
	if len(user) < 1 {
		tx := DB.Create(&User{
			Admin:      true,
			Name:       "admin",
			LastName:   "admin",
			MiddleName: "admin",
			Email:      "admin@admin.ru",
			Password:   fmt.Sprintf("%x", sha256.Sum256([]byte("admin"))),
		})
		if tx.Error != nil {
			logging.Print.Error(tx.Error)
			return
		}

		if err := os.MkdirAll(fmt.Sprintf("lib/users/%v", strconv.Itoa(int(tx.Statement.Model.(*User).Model.ID))), 0777); err != nil {
			logging.Print.Error(err)
			return
		}

	}
	logging.Print.Infof("Available user admin")
	fmt.Println("                       login: admin@admin.ru")
	fmt.Println("                       password: admin")
}

func insertRoles() {
	var roles []Role
	DB.Find(&roles)
	if len(roles) == 0 {
		data := []string{
			"Android разработчик",
			"iOS разработчик",
			"Frontend разработчик",
			"Backend разработчик",
			"Аналитик",
			"Тестировщик",
			"Администратор",
			"Менеджер",
			"Дизайнер",
		}
		for _, v := range data {
			DB.Create(&Role{
				Name: v,
			})
		}
	}
}
