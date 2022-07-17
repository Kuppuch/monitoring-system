package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type User struct {
	gorm.Model
	Admin      bool   `json:"admin"`
	Name       string `json:"name"`
	LastName   string `json:"lastname"`
	MiddleName string `json:"middlename"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
}

func GetAllUsers() []User {
	var users []User
	tx := DB.Find(&users)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return users
}

func (u *User) GetUser() error {
	tx := DB.Where("id = ?", u.ID).Find(u)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u *User) GetUserByEmail() error {
	tx := DB.Where("email = ?", u.Email).Find(u)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u User) InsertUser() int64 {
	tx := DB.Create(&u)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}
