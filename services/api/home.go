package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func getMainPage(c *gin.Context) {
	bts, err := ioutil.ReadFile("./pages/img/logo.png")
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string
	mimeType := http.DetectContentType(bts)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(bts)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"image": base64Encoding,
	})
}
