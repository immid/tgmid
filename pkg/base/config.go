package base

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	ConfigDir    = "./configs/"
	config       *viper.Viper
	configLoaded = false
)

func GetConfig() *viper.Viper {
	if configLoaded {
		return config
	}
	viper.AddConfigPath(ConfigDir)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("telegram config: %s\n", err))
	}
	config = viper.GetViper()
	configLoaded = true
	filename := config.GetString("tdlib.log_file")
	if len(filename) == 0 || filename == "stderr" {
		log.SetOutput(os.Stderr)
	} else if filename == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		filename = config.GetString("tdlib.data_dir") + "/" + filename
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err == nil {
			log.SetOutput(logFile)
		}
	}
	return config
}
