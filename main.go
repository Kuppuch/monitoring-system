package main

import (
	_ "monitoring-system/config"
	"monitoring-system/services/api"
	"monitoring-system/services/logging"
	"monitoring-system/services/middleware"
)

func main() {
	logging.Init()
	err := middleware.Connector()
	if err != nil {
		logging.Print.Fatal(err)
	}
	logging.Print.Info("Listen on port 25595 http://localhost:25595/")
	api.Router()
}
