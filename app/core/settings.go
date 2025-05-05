package core

import "log"
import "github.com/spf13/viper"

type Webhook struct {
	TokenRevoked string `mapstructure:"token_revoked"`
}

type ListServices struct {
	Service string            `mapstructure:"service"`
	Webhook Webhook           `mapstructure:"webhook"`
	Auth    []string          `mapstructure:"auth"`
	Role    map[string]string `mapstructure:"role"`
}

type ServiceConfig struct {
	AppIp      string `mapstructure:"APP_IP"`
	AppEnv     string `mapstructure:"APP_ENV"`
	PgUsername string `mapstructure:"POSTGRES_USER"`
	PgPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PgHost     string `mapstructure:"POSTGRES_HOST"`
	PgPort     string `mapstructure:"POSTGRES_PORT"`
	PgDatabase string `mapstructure:"POSTGRES_DB"`

	MoUsername string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MoPassword string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	MoHost     string `mapstructure:"MONGO_HOST"`
	MoPort     string `mapstructure:"MONGO_PORT"`
	MoDatabase string `mapstructure:"MONGO_INITDB_DATABASE"`

	AccessTokenSecret  string                  `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string                  `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenMinute  int                     `mapstructure:"ACCESS_TOKEN_MINUTE"`
	RefreshTokenMinute int                     `mapstructure:"REFRESH_TOKEN_MINUTE"`
	ListServices       map[string]ListServices `mapstructure:"list_services"`
	NameApp            string                  `mapstructure:"name_app"`
	ShortJwtValue      string                  `mapstructure:"short_jwt_value"`
}

func (s *ServiceConfig) ReadConfig(pathEnv string, pathYaml string) {
	// Чтение конфигурации из ENV
	viper.SetConfigFile(pathEnv)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("ENV ERROR = ", err)
	}

	// Чтение YAML конфигурации и объединение с ENV
	viper.SetConfigFile(pathYaml)
	err = viper.MergeInConfig()
	if err != nil {
		log.Fatal("YAML ERROR = ", err)
	}

	// Выполняем Unmarshal один раз после объединения
	err = viper.Unmarshal(s)
	if err != nil {
		log.Fatal("Конфигурация не может быть загружена: ", err)
	}
}
