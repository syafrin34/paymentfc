package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() Config {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./files/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error with config file %v", err)

	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("error unmarshall config: %v", err)
	}
	tempDebug31, _ := json.Marshal(cfg)
	fmt.Printf("\n======= Debug config.go - line 31 ==== \n\n%s\n\n===============\n\n\n, ", string(tempDebug31))

	return cfg
}
