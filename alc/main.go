package main

import (
	"fmt"
	"github.com/GreatGodApollo/alc/prompt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("alc")
	viper.SetConfigType("json")

	// Api Settings
	viper.SetDefault("api.baseurl", "http://localhost:8080")

	// Auth Settings
	viper.SetDefault("auth.username", "username")
	viper.SetDefault("auth.password", "password")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file", viper.ConfigFileUsed())
	} else {
		panic("Could not load configuration file: " + err.Error())
	}
}

func main() {
	restyCli := resty.New()

	prompt.RunPrompt(restyCli)
}
