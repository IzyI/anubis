package core

import "log"
import "github.com/spf13/viper"

type ServiceConfig struct {
	AppIp              string `mapstructure:"APP_IP"`
	AppEnv             string `mapstructure:"APP_ENV"`
	PgUsername         string `mapstructure:"POSTGRES_USER"`
	PgPassword         string `mapstructure:"POSTGRES_PASSWORD"`
	PgHost             string `mapstructure:"POSTGRES_HOST"`
	PgPort             string `mapstructure:"POSTGRES_PORT"`
	PgDatabase         string `mapstructure:"POSTGRES_DB"`
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenHour    int    `mapstructure:"ACCESS_TOKEN_HOUR"`
	RefreshTokenHour   int    `mapstructure:"REFRESH_TOKEN_HOUR"`
	Server             struct {
		Phone []string `mapstructure:"phone"`
		OAuth []string `mapstructure:"oauth"`
	} `mapstructure:"server"`
}

func (s *ServiceConfig) ReadConfig(pathEnv string, pathyaml string) {
	viper.SetConfigFile(pathEnv)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("ENV ERROR = ", err)
	}

	err = viper.Unmarshal(&s)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	viper.SetConfigFile(pathyaml)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("YAML ERROR = ", err)
	}

	err = viper.Unmarshal(&s)
	if err != nil {
		log.Fatal("Yaml can't be loaded: ", err)
	}

}
