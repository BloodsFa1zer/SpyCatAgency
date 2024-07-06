package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	UserDBPassword string `env:"POSTGRES_PASSWORD"`
	UserDBName     string `env:"POSTGRES_USER"`
	DBName         string `env:"POSTGRES_DB"`
	DriverDBName   string `env:"POSTGRES_DRIVER"`
	DBPort         string `env:"POSTGRES_PORT"`

	DBHost       string `env:"DB_HOST"`
	PostgresHost string `env:"POSTGRES_HOST"`
}

var cfg *Config

func LoadENV(filename string) *Config {
	if cfg != nil {
		return cfg
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Panic().Err(err).Msg("Error loading .env file")
	}
	log.Info().Msg("Successfully loaded .env")

	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic().Err(err).Msg("Error parsing environment variables")
	}
	log.Info().Msg("Successfully parsed .env")

	return cfg
}
