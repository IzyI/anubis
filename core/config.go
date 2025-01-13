package core

type Env struct {
	AppIp              string `mapstructure:"APP_IP"`
	AppEnv             string `mapstructure:"APP_ENV"`
	PgUsername         string `mapstructure:"PG_USERNAME"`
	PgPassword         string `mapstructure:"PG_PASSWORD"`
	PgHost             string `mapstructure:"PG_HOST"`
	PgPort             string `mapstructure:"PG_PORT"`
	PgDatabase         string `mapstructure:"PG_DATABASE"`
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenHour    int    `mapstructure:"ACCESS_TOKEN_HOUR"`
	RefreshTokenHour   int    `mapstructure:"REFRESH_TOKEN_HOUR"`
}
