package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigModel struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

// Can be used to access configuration from other packeages
var Config ConfigModel

// Loads configuration from Environment and .env file
// .env overwrites the Environment
func LoadConfig() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("no .env Found. Using only Variables from the Environment")
	}

	viper.Unmarshal(&Config)
}
