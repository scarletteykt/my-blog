package config

import (
	"github.com/joho/godotenv"
	"github.com/scraletteykt/my-blog/pkg/logger"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
)

type (
	Config struct {
		Auth     AuthConfig
		HTTP     HTTPConfig
		Postgres PostgresConfig
	}

	AuthConfig struct {
		Secret string
	}

	HTTPConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
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

func NewConfig(log logger.Logger) (*Config, error) {
	populateDefaults()

	if err := parseConfig(log); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.Auth.Secret = os.Getenv("SECRET_KEY")
	cfg.Postgres.Password = os.Getenv("DB_PASSWORD")
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

func parseConfig(log logger.Logger) error {
	viper.AddConfigPath(".")
	viper.SetConfigFile("configs/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
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
	viper.SetDefault("http.maxHeaderBytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}
