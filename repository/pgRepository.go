package repository

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PgRepository base struct from a postgres repository
// contains the database config after a
// successful connection
type PgRepository struct {
	DB *gorm.DB
}

// PgRepositoryConnConfig struct to config the connection to postgres
type PgRepositoryConnConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	TimeZone     string
}

func NewPgRepository(pgf PgRepositoryConnConfig) (*PgRepository, error) {

	var pgDNS string = fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		pgf.Username,
		pgf.Password,
		pgf.Host,
		pgf.Port,
		pgf.DatabaseName)

	var repo PgRepository
	//var err error

	//TODO: Implement interface to use zerolog in gorm package
	db, err := gorm.Open(postgres.Open(pgDNS), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("error while trying to connect to postgres database")
		return nil, err
	}

	repo.DB = db

	return &repo, nil

}
