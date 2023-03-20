package main

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"monitoring-system/notify_service/api"
	"monitoring-system/notify_service/pkg"
	"monitoring-system/services/logging"
)

func main() {
	err := pkg.Connector()
	if err != nil {
		logging.Print.Error(err)
	}

	r := gin.Default()
	r.Use(cors.AllowAll())
	r.POST("/email", api.SendEmail)
	r.POST("/email_code", api.GetCodeByEmail)
	err = r.Run(":25596")
	if err != nil {
		logging.Print.Error(err)
	}
}
