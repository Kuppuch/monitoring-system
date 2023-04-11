package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"monitoring-system/services/api/socket"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
	"net/http"
)

var WS *websocket.Conn
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		//return "chrome-extension://fgponpodhbmadfljofbimhhlengambbn" == "chrome-extension://fgponpodhbmadfljofbimhhlengambbn"
		return true
	},
}

var mapPool = make(map[string]*socket.Pool)
var Pool *socket.Pool

func init() {
	Pool = socket.NewPool()
	go Pool.Run()
}

func sendMessage(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		logging.Print.Error(err)
		c.JSON(http.StatusBadRequest, gin.Error{
			Err:  err,
			Type: 0,
			Meta: "error by get raw data",
		})
		return
	}

	m := map[string]string{}
	err = json.Unmarshal(raw, &m)
	socket.BigChannel <- []byte(m["mess"])
}

func socketFunc(c *gin.Context) {
	socket.ServeWs(Pool, c.Writer, c.Request)
}

func getNotificationPage(c *gin.Context) {
	user, _ := GetUserByToken(c)
	notifications := middleware.GetAssignedToNotification(user.ID)
	c.HTML(http.StatusOK, "notification.html", gin.H{"user": user, "notifications": notifications})
}

func getNotifications(c *gin.Context) {
	user, _ := GetUserByToken(c)
	notifications := middleware.GetUnreadNotification(user.ID)
	c.JSON(http.StatusOK, notifications)
}
