package core

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
