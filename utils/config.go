package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	FolderPath     string `mapstructure:"folderPath"`
	OutPutFileName string `mapstructure:"outPutFileName"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path) // <- to work with Dockerfile setup
	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
