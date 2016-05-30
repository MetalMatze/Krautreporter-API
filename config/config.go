package config

import (
	"github.com/MetalMatze/gollection"
)

// GetConfig returns the gollection configuration
func GetConfig() gollection.Config {
	gollection.LoadEnv(".env") // Loads env variables from .env, but doesn't override if already set

	return gollection.Config{
		AppConfig: gollection.AppConfig{
			Name:  "Kraureporter-API",
			Usage: "RESTful json API which crawls krautreporter.de and serves its content",
			Host:  gollection.GetEnv("APP_HOST", "127.0.0.1"),
			Port:  gollection.GetEnvInt("APP_PORT", 1234),
			Debug: gollection.GetEnvBool("APP_DEBUG", false),
		},
		DBConfig: gollection.DBConfig{
			Dialect:  gollection.GetEnv("DB_DIALECT", "postgres"), // postgres, mysql or sqlite3
			Host:     gollection.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     gollection.GetEnvInt("DB_PORT", 5432),
			Database: gollection.GetEnv("DB_DATABASE", "postgres"),
			Username: gollection.GetEnv("DB_USERNAME", "postgres"),
			Password: gollection.GetEnv("DB_PASSWORD", "postgres"),
		},
		RedisConfig: gollection.RedisConfig{
			Host:     gollection.GetEnv("REDIS_HOST", "127.0.0.1"),
			Password: gollection.GetEnv("REDIS_PASSWORD", ""),
			Port:     gollection.GetEnvInt("REDIS_PORT", 6379),
			Database: 0,
		},
	}
}
