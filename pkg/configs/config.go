package configs

import (
	"log"

	"zocket/pkg/utils"

	"github.com/spf13/viper"
)

func GetConfig() (*Config, error) {

	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file:", err)
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Println("Unable to decode into struct:", err)
		return nil, err
	}

	if err := utils.NewValidator().Struct(config); err != nil {
		log.Println("Config validation failed:", err)
		return nil, err
	}

	return &config, nil
}
