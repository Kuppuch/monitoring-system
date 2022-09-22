package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var (
	DbName    string
	DbHost    string
	DBPort    string
	DbSslMode string

	DbUser     string
	DbPassword string
)

func init() {
	configFile := flag.String("config", "config/config.json", "config which should be used")
	flag.Parse()

	viper.SetConfigType("json")
	viper.SetConfigFile(*configFile)

	if err := viper.ReadInConfig(); err != nil {
		*configFile = "../config/config.json"
		viper.SetConfigFile(*configFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("Unable to read config file: %s", err))
		}
	}

	DbUser = viper.GetString("db.user")
	DbName = viper.GetString("db.dbname")
	DbHost = viper.GetString("db.host")
	DBPort = viper.GetString("db.port")
	DbSslMode = viper.GetString("db.sslmode")
	DbPassword = viper.GetString("db.password")

	createLibRep()
}

func createLibRep() {
	err := os.MkdirAll("lib/users", 0777)
	if err != nil {
		panic(fmt.Sprintf("Unable to create lib/user directory: %s", err))
	}
}
