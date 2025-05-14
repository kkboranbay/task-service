package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

type AuthConfig struct {
	JWTSecret        string
	TokenExpireDelta time.Duration
}

type LoggerConfig struct {
	Level string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "5s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "5s")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "taskdb")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_CONNS", 10)
	viper.SetDefault("DB_TIMEOUT", "5s")

	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRE_DELTA", "24h")

	viper.SetDefault("LOG_LEVEL", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
		}
	}

	var config Config

	readTimeout, err := time.ParseDuration(viper.GetString("SERVER_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга SERVER_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := time.ParseDuration(viper.GetString("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга SERVER_WRITE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := time.ParseDuration(viper.GetString("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}

	config.Server = ServerConfig{
		Port:            viper.GetString("SERVER_PORT"),
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
	}

	dbTimeout, err := time.ParseDuration(viper.GetString("DB_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DB_TIMEOUT: %w", err)
	}

	config.Database = DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
		MaxConns: viper.GetInt("DB_MAX_CONNS"),
		Timeout:  dbTimeout,
	}

	tokenExpireDelta, err := time.ParseDuration(viper.GetString("JWT_EXPIRE_DELTA"))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JWT_EXPIRE_DELTA: %w", err)
	}

	config.Auth = AuthConfig{
		JWTSecret:        viper.GetString("JWT_SECRET"),
		TokenExpireDelta: tokenExpireDelta,
	}

	config.Logger = LoggerConfig{
		Level: viper.GetString("LOG_LEVEL"),
	}

	return &config, nil
}
