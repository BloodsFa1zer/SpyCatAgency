package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"spyCat/config"
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase() *Database {
	cfg := config.LoadENV(".env")

	connStr := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBPort, cfg.UserDBName, cfg.UserDBPassword, cfg.DBName)
	// for docker file
	//connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	//	cfg.PostgresHost, cfg.DBPort, cfg.UserDBName, cfg.UserDBPassword, cfg.DBName)

	log.Info().Msgf("Connection string: %s", connStr)

	db, err := sql.Open(cfg.DriverDBName, connStr)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to connect to database")
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to ping the database")
		return nil
	}
	log.Info().Msg("Successfully connected to the database.")

	return &Database{Connection: db}
}
