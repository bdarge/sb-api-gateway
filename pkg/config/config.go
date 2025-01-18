package config

import "github.com/spf13/viper"
import 	"golang.org/x/exp/slog"

// Config holds config for this repo
type Config struct {
	Port              string 		 `mapstructure:"PORT"`
	AuthSvcURL        string 		 `mapstructure:"AUTH_SVC_URL"`
	APISvcURL         string 		 `mapstructure:"API_SVC_URL"`
	BaseURL           string 		 `mapstructure:"BASE_URL"`
	RefreshTokenExpOn int    		 `mapstructure:"JWT_REFRESH_TOKEN_EXP_ON"`
	UIDomain          string 		 `mapstructure:"UI_DOMAIN"`
	CurrencySvcURL    string 		 `mapstructure:"CURRENCY_SVC_URL"`
	LogLevel 					slog.Level `mapstructure:"LOG_LEVEL"`
}

// LoadConfig load config
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
