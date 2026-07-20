package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type ConfigProvider interface {
	GetAddress() string
	GetDBConfig() *DBConfig
}

type Config struct {
	AppHost  string `mapstructure:"HOST"`
	AppPort  string `mapstructure:"PORT"`
	DbConfig *DBConfig
}

func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")

	var config Config
	var dbConfig DBConfig

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("error reading config file, %s", err)
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config, %v", err)
	}
	if err := viper.Unmarshal(&dbConfig); err != nil {
		return nil, fmt.Errorf("unable to decode dbConfig, %v", err)
	}

	if config.AppHost == "" {
		config.AppHost = "localhost"
	}
	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	config.DbConfig = &dbConfig

	return &config, nil
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.AppHost, c.AppPort)
}

func (c *Config) GetDBConfig() *DBConfig {
	return c.DbConfig
}

type DBConfig struct {
	Host string `mapstructure:"DB_HOST"`
	Port string `mapstructure:"DB_PORT"`
	User string `mapstructure:"DB_USER"`
	Pass string `mapstructure:"DB_PASSWORD"`
	Name string `mapstructure:"DB_NAME"`
}

func (dbConfig *DBConfig) GetConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
}
