package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
	Upstream UpstreamConfig `mapstructure:",squash"`
	Logger   LoggerConfig   `mapstructure:",squash"`
}

type AppConfig struct {
	Name     string `mapstructure:"app_name"`
	Env      string `mapstructure:"app_env"`
	HTTPPort int    `mapstructure:"app_http_port"`
}

type JWTConfig struct {
	AccessSecret string `mapstructure:"jwt_access_secret"`
}

type UpstreamConfig struct {
	IdentityServiceURL   string `mapstructure:"identity_service_url"`
	RestaurantServiceURL string `mapstructure:"restaurant_service_url"`
}

type LoggerConfig struct {
	Level string `mapstructure:"logger_level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	// Set defaults
	viper.SetDefault("app_name", "gopher-gateway-service")
	viper.SetDefault("app_env", "development")
	viper.SetDefault("app_http_port", 8000)
	viper.SetDefault("logger_level", "debug")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
