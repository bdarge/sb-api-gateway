package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	AuthSvcUrl        string `mapstructure:"AUTH_SVC_URL"`
	ApiSvcUrl         string `mapstructure:"API_SVC_URL"`
	BaseUrl           string `mapstructure:"BASE_URL"`
	RefreshTokenExpOn int    `mapstructure:"JWT_REFRESH_TOKEN_EXP_ON"`
	UIDomain          string `mapstructure:"UI_DOMAIN"`
}

func LoadConfig(target string) (c Config, err error) {
	viper.AddConfigPath("./envs")
	viper.SetConfigName(target)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
