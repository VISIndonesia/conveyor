package config

import (
	"log"

	"github.com/spf13/viper"
)

type PubSub struct {
	CredFile   string
	ProjectID  string
	Subscriber string
	Timeout    int
}

type BigQuery struct {
	CredFile  string
	ProjectID string
	Dataset   string
}

type Config struct {
	PubSub   PubSub
	BigQuery BigQuery
	Gap      int
}

var Configuration Config

func init() {
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.conveyor")
	viper.AddConfigPath("/mnt")
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	viper.SetDefault("gap", 10)

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
