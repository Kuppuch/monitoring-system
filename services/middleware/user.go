package middleware

import (
	"gorm.io/gorm"
	"monitoring-system/services/logging"
)

type User struct {
	gorm.Model
	Admin       bool        `json:"admin"`
	Name        string      `json:"name"`
	LastName    string      `json:"lastname"`
	MiddleName  string      `json:"middlename"`
	Email       string      `json:"email"`
	Password    string      `json:"password,omitempty"`
	CreateIssue []Issue     `gorm:"foreignKey:CreatorID" json:"create_issue,omitempty"`
	AssignIssue []Issue     `gorm:"foreignKey:AssignedToID" json:"assign_issue,omitempty"`
	Auth        []Auth      `gorm:"foreignKey:UserID" json:"auth,omitempty"`
	Members     []Member    `gorm:"foreignKey:UserID" json:"members,omitempty"`
	Timespent   []Timespent `gorm:"foreignKey:UserID" json:"timespent,omitempty"`
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

func (u *User) Insert() int64 {
	tx := DB.Create(u)
	if tx.Error != nil {
		logging.Print.Warning(tx.Error)
	}
	return tx.RowsAffected
}

func (u *User) Update() {
	user := User{}
	DB.Where("id = ?", u.ID).Find(&user)
	u.Password = user.Password
	DB.Model(u).Save(&u)
}
