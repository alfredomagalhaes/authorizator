package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgConnConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	TimeZone     string
}

// Creates the connection to the database
// using gorm lib to be used all over the application
func NewPgConnection(config PgConnConfig) (*gorm.DB, error) {
	var pgDNS string = fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName)

	//TODO: Implement interface to use zerolog in gorm package
	db, err := gorm.Open(postgres.Open(pgDNS), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("error while trying to connect to postgres database")
		return nil, err
	}

	return db, nil
}
