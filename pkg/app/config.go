package app

import "github.com/spf13/viper"

func InitConfig(fileName string) error {
	viper.AddConfigPath("configs")
	viper.SetConfigName(fileName)
	return viper.ReadInConfig()
}
