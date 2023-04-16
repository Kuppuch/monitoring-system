package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
	"os"
	"strconv"
)

func getUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.DefaultQuery("uid", "0"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse int",
		})
	}
	if uid == 0 {
		c.HTML(http.StatusOK, "users.html", gin.H{"users": middleware.GetAllUsers()})
		//c.JSON(http.StatusOK, middleware.GetAllUsers())
		return
	}
	user := middleware.User{}
	user.ID = uint(uid)
	err = user.GetUser()
	if err != nil {
		logging.Print.Warning(err)
	}
	user.Password = ""
	//c.JSON(http.StatusOK, user)

	//Profile photo load
	bts, err := ioutil.ReadFile("./lib/users/" + strconv.Itoa(int(user.ID)) + "/photo.png")
	if err != nil {
		log.Fatal(err)
	}
	base64Encoding := "data:image/png;base64," + base64.StdEncoding.EncodeToString(bts)

	c.HTML(http.StatusOK, "users.html", gin.H{
		"user":  user,
		"image": base64Encoding,
	})
}

// начало админских функций по пользователям

func getPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func insertUser(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		logging.Print.Warning(err)
	}
	user := middleware.User{}
	err = json.Unmarshal(raw, &user)
	if err != nil {
		logging.Print.Warning(err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "perhaps raw is null",
		})
		return
	}
	//TODO переделать на длину
	if len(user.Email) < 1 || (len(user.Password) < 1 && user.ID == 0) || len(user.LastName) < 1 || len(user.Name) < 1 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("null field"),
			Type: 0,
			Meta: "not enough data to register",
		})
		return
	}

	dbUser := middleware.User{}
	dbUser.Email = user.Email
	err = dbUser.GetUserByEmail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed check existing user",
		})
		logging.Print.Error(err)
		return
	}
	if dbUser.ID > 0 {
		user.Update()
		c.JSON(http.StatusOK, middleware.GetSuccess())
		return
	}

	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	if err != nil {
		logging.Print.Warning(err)
		c.JSON(http.StatusInternalServerError, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed create password hash",
		})
		return
	}

	rowAffected := user.Insert()
	if rowAffected > 0 {
		err = createUserRep(user.ID)
		if err != nil {
			logging.Print.Error("Unable to create user directory: ", err)
		}
		c.JSON(http.StatusOK, middleware.GetSuccess())
	} else {
		c.JSON(http.StatusBadRequest, middleware.GetBadRequest())
	}

}

func getUserAdm(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "failed parse int",
		})
	}
	if uid == 0 {
		c.HTML(http.StatusOK, "users.html", gin.H{"users": middleware.GetAllUsers()})
		//c.JSON(http.StatusOK, middleware.GetAllUsers())
		return
	}
	user := middleware.User{}
	user.ID = uint(uid)
	err = user.GetUser()
	if err != nil {
		logging.Print.Warning(err)
	}
	user.Password = ""
	c.HTML(http.StatusOK, "register.html", gin.H{"user": user})
}

// конец админских функций

func createUserRep(userID uint) error {
	err := os.Mkdir("lib/users/"+strconv.Itoa(int(userID)), 0777)
	if err != nil {
		return err
	}
	return nil
}

func uploadProfileImg(c *gin.Context) {
	src, _ := c.FormFile("image")
	file, _ := src.Open()
	defer file.Close()

	if src.Size > 1000001 {
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("too big file"),
			Type: 0,
			Meta: "file must not exceed 20 mb",
		})
	}

	var bb bytes.Buffer
	imageBytes := make([]byte, src.Size+128)
	_, err := file.Read(imageBytes)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	bb.Write(imageBytes)

	dataType := http.DetectContentType(imageBytes)
	switch dataType {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			fmt.Println(err, "unable to decode jpeg")
		}

		if err := png.Encode(&bb, img); err != nil {
			fmt.Println(err, "unable to encode png")
		}
	case "image/png":
		bb.Write(imageBytes)
	default:
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  errors.New("uploading failed"),
			Type: 0,
			Meta: "unsupported format or file is not image",
		})
	}
	user, _ := GetUserByToken(c)
	profile, _ := os.Create("lib/users/" + strconv.Itoa(int(user.ID)) + "/photo.png")
	defer profile.Close()
	profile.Write(bb.Bytes())
}

func GetProfilePhoto(c *gin.Context) {
	id := c.Param("id")
	_, err := ioutil.ReadFile("./lib/users/" + id + "/photo.png")
	if err != nil {
		logging.Print.Warning(err)
	}
	user, _ := GetUserByToken(c)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": user,
	})
}
