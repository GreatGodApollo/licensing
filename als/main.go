package main

import (
	"database/sql"
	"fmt"
	"github.com/GreatGodApollo/als/database"
	"github.com/GreatGodApollo/als/server"
	"github.com/GreatGodApollo/als/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("als")
	viper.SetConfigType("json")

	// Server Defaults
	viper.SetDefault("server.bind", ":8080")
	viper.SetDefault("server.production", false)

	// Database Defaults
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "root")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "3306")
	viper.SetDefault("db.name", "license")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file", viper.ConfigFileUsed())
	} else {
		panic("Could not load configuration file: " + err.Error())
	}

	// Cryptography Key
	if viper.GetString("crypt.key") == "" {
		viper.Set("crypt.key", utils.RandomString(16))
		if err := viper.WriteConfig(); err != nil {
			panic("Could not generate crypto key: " + err.Error())
		}
	}
}

func main() {
	// Database Setup
	db, err := database.Setup()
	if err != nil {
		panic("Could not set up database: " + err.Error())
	}
	defer db.Close()

	server.Setup(db)
	server.RunAPI()
}
