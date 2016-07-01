package config

import (
	"log"
	"os"
	"strconv"

	"github.com/gollection/gollection"
	"github.com/gollection/gollection/database"
	"github.com/gollection/gollection/router"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Config         gollection.Config
	DatabaseConfig database.Config
	RouterConfig   router.Config
}

func Config() AppConfig {
	loadEnv()

	databasePort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Println(err)
	}

	routerPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Println(err)
	}

	return AppConfig{
		Config: gollection.Config{
			Name:  "Kraureporter-API",
			Usage: "RESTful json API which crawls krautreporter.de and serves its content",
		},
		// Use all for non sqlite databases
		DatabaseConfig: database.Config{
			Driver:   "postgres",
			Host:     os.Getenv("DB_HOST"),
			Port:     databasePort,
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Database: os.Getenv("DB_NAME"),
		},
		RouterConfig: router.Config{
			Host: os.Getenv("APP_HOST"),
			Port: routerPort,
		},
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}
