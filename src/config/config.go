package config

import (
	"fmt"
	env "github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type flatEnv struct {
	AppName string `env:"APP_NAME"`
	AppEnv  string `env:"APP_ENV"` // local|dev|uat|sit|prod

	LoggerLevel  string `env:"LOGGER_LEVEL"`
	LoggerOutput string `env:"LOGGER_OUTPUT"`

	GRPCPort int `env:"GRPC_PORT"`

	GRPCRestProxyPort int `env:"GRPC_REST_PROXY_PORT"`

	PostgresHost        string `env:"POSTGRES_HOST"`
	PostgresPort        int    `env:"POSTGRES_PORT"`
	PostgresDatabase    string `env:"POSTGRES_DATABASE"`
	PostgresUsername    string `env:"POSTGRES_USERNAME"`
	PostgresPassword    string `env:"POSTGRES_PASSWORD,unset"`
	PostgresLogging     bool   `env:"POSTGRES_LOGGING"`
	PostgresConnMaxOpen int    `env:"POSTGRES_CONN_MAX_OPEN"`
	PostgresConnMaxIdle int    `env:"POSTGRES_CONN_MAX_IDLE"`
	PostgresSSLMode     string `env:"POSTGRES_SSL_MODE"`
}

type AppInfo struct {
	Name        string
	Environment string
}

type LoggerConfig struct {
	Level string
}

type GRPCServerConfig struct {
	Port int
}

type GRPCRESTProxyConfig struct {
	Port int
}

type PostgreSQLConfig struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	Logging      bool
}

type Config struct {
	AppInfo             AppInfo
	Logger              LoggerConfig
	GRPCServer          GRPCServerConfig
	GRPCRESTProxyServer GRPCRESTProxyConfig
	PprofEnable         bool
	PostgreSQL          PostgreSQLConfig
}

func LoadFromEnv() (*Config, error) {
	err := godotenv.Load() // load dari .env
	if err != nil {
		fmt.Println("warning: no .env file found, continuing...")
	}
	var envCfg flatEnv
	err = env.Parse(&envCfg)
	if err != nil {
		return nil, err
	}
	// ========== BASE CONFIG ==========
	cfg := &Config{
		AppInfo: AppInfo{
			Name:        envCfg.AppName,
			Environment: envCfg.AppEnv,
		},
		Logger: LoggerConfig{
			Level: envCfg.LoggerLevel,
		},
		GRPCServer: GRPCServerConfig{
			Port: envCfg.GRPCPort,
		},
		GRPCRESTProxyServer: GRPCRESTProxyConfig{
			Port: envCfg.GRPCRestProxyPort,
		},
		PprofEnable: false,
		PostgreSQL: PostgreSQLConfig{
			Host:         envCfg.PostgresHost,
			Port:         envCfg.PostgresPort,
			Database:     envCfg.PostgresDatabase,
			Username:     envCfg.PostgresUsername,
			Password:     envCfg.PostgresPassword,
			SSLMode:      envCfg.PostgresSSLMode,
			MaxOpenConns: envCfg.PostgresConnMaxOpen,
			MaxIdleConns: envCfg.PostgresConnMaxIdle,
			Logging:      false,
		},
	}
	return cfg, nil
}
