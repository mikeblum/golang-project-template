package conf

import (
	"os"

	"github.com/spf13/viper"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

func NewConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	config.AddConfigPath("..")
	config.AddConfigPath(GetEnv(EnvConfigPath, "/"))
	config.AutomaticEnv()
	err := config.ReadInConfig()
	return config, err
}

// GetEnv lookup an environment variable or fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
