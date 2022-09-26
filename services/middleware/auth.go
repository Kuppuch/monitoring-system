package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	ID uint
}

type Auth struct {
	gorm.Model
	UserID uint
	Token  string
}

func (a Auth) InsertAuth() error {
	tx := DB.Create(&a)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetAuthByUserID(userID uint) Auth {
	a := Auth{}
	DB.Where("user_id = ?", userID).Order("created_at desc").Limit(1).Find(&a)
	return a
}

func GetToken(id uint) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(), //10 часов
			IssuedAt:  time.Now().Unix(),
		},
		ID: id,
	})
	token, _ := t.SignedString([]byte("123"))
	return token
}

func CheckToken(accesstoken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("123"), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	user := User{
		Model: gorm.Model{ID: claims.ID},
	}
	err = user.GetUser()
	if ok && token.Valid && user.ID > 0 {
		return claims.ID, nil
	}
	return 0, errors.New("invalid access token")
}
