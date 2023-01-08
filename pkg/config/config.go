package config

import "github.com/spf13/viper"

type Config struct {
	Port       string `mapstructure:"PORT"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
	ApiSvcUrl  string `mapstructure:"API_SVC_URL"`
	// base url
	BaseUrl string `mapstructure:"BASE_URL"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
