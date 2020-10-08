package config

import "github.com/spf13/viper"

var (
	ConfigHandler = viper.New()
)

func init() {
	ConfigHandler.SetConfigName("config")
	ConfigHandler.SetConfigType("yaml")
	ConfigHandler.AddConfigPath(".")
	ConfigHandler.ReadInConfig()
}
