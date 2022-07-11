package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	defaultHTTPPort = "8080"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Auth     AuthConfig
	}

	AuthConfig struct {
		Secret string
	}

	HTTPConfig struct {
		Port string `mapstructure:"port"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `sslmode:"host"`
	}
)

func InitConfig() (*Config, error) {
	populateDefaults()

	if err := parseConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.Postgres.Password = os.Getenv("DB_PASSWORD")
	cfg.Auth.Secret = os.Getenv("SECRET_KEY")
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	return nil
}

func parseConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile("configs/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return err
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
		return err
	}
	return nil
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
}
