package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"monitoring-system/services/logging"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		//return "chrome-extension://fgponpodhbmadfljofbimhhlengambbn" == "chrome-extension://fgponpodhbmadfljofbimhhlengambbn"
		return true
	},
}

var Chanel = make(chan string)

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
	Chanel <- m["mess"]
}

func socket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		//Read Message from client
		//mt, message, err := ws.ReadMessage()
		//if err != nil {
		//	fmt.Println(err)
		//	break
		//}
		////If client message is ping will return pong
		//if string(message) == "ping" {
		//	message = []byte("pong")
		//}
		////Response message to client
		//err = ws.WriteMessage(mt, message)
		//if err != nil {
		//	fmt.Println(err)
		//	break
		//}

		select {
		case info := <-Chanel:
			err = ws.WriteMessage(1, []byte(info))
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		//for i := 0; i < 5; i++ {
		//	err = ws.WriteMessage(1, []byte(strconv.Itoa(i)))
		//	if err != nil {
		//		fmt.Println(err)
		//		break
		//	}
		//	time.Sleep(time.Second)
		//}
	}
}
