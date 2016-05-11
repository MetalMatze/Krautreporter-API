package main

import (
	"github.com/MetalMatze/gollection"
)

func GetConfig() gollection.Config {
	gollection.LoadEnv(".env") // Loads env variables from .env, but doesn't override if already set

	return gollection.Config{
		AppConfig: gollection.AppConfig{
			Name:  "Kraureporter-API",
			Usage: "RESTful json API which crawls krautreporter.de and serves its content",
			Host:  gollection.GetEnv("APP_HOST", "localhost"),
			Port:  gollection.GetEnvInt("APP_PORT", 1234),
			Env:   gollection.GetEnv("APP_ENV", "production"), // local, testing, production
		},
		DBConfig: gollection.DBConfig{
			Dialect:  gollection.GetEnv("DB_DIALECT", "postgres"), // postgres, mysql or sqlite3
			Host:     gollection.GetEnv("DB_HOST", "localhost"),
			Port:     gollection.GetEnvInt("DB_PORT", 5432),
			Database: gollection.GetEnv("DB_DATABASE", "postgres"),
			Username: gollection.GetEnv("DB_USERNAME", "postgres"),
			Password: gollection.GetEnv("DB_PASSWORD", "postgres"),
		},
	}
}
