package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getMainPage(c *gin.Context) {
	file, err := os.Open("./pages/img/logo.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//bytess, err := ioutil.ReadFile("./pages/img/logo.png")
	bytess, err := ioutil.ReadFile("logo1.png")
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytess)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += string(bytess)

	// Print the full base64 representation of the image
	fmt.Println(base64Encoding)

	//src, _ := png.Decode(file)
	//buffer := bytes.Buffer{}
	//_ = png.Encode(&buffer, src)
	//imgBase64 := base64.StdEncoding.EncodeToString([]byte(buffer))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"image": base64Encoding,
	})
}
