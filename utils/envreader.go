package utils

import "github.com/spf13/viper"

func Init() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
