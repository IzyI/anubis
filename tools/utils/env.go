package utils

import "log"
import "github.com/spf13/viper"

func ReadConfig(config interface{}) {
	viper.SetConfigFile("config/.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("ENV ERROR = ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	viper.SetConfigFile("config/server.yaml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("YAML ERROR = ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Yaml can't be loaded: ", err)
	}

}
