package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
)

type (
	Config struct {
		HTTP     HTTPConfig
		Postgres PostgresConfig
	}

	HTTPConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"read_timeout"`
		WriteTimeout       time.Duration `mapstructure:"write_timeout"`
		MaxHeaderMegabytes int           `mapstructure:"max_header_megabytes"`
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

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
}

func parseConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigFile("config")
	return viper.ReadInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.read_timeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.write_timeout", defaultHTTPRWTimeout)
}
