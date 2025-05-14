package provider

import "github.com/spf13/viper"

func LoadViperConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$PWD/.config")
	if err := v.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	return v
}
