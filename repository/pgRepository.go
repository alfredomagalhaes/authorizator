package repository

import (
	"fmt"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PgRepository base struct from a postgres repository
// contains the database config after a
// successful connection
type PgRepository struct {
	DB  *gorm.DB
	log *zerolog.Logger
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

func NewPgRepository(pgf PgRepositoryConnConfig, log *zerolog.Logger) (*PgRepository, error) {
	var pgDNS string = fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s`,
		pgf.Host,
		pgf.Username,
		pgf.Password,
		pgf.DatabaseName,
		pgf.Port,
		pgf.TimeZone)
	var repo *PgRepository
	var err error

	repo.DB, err = gorm.Open(postgres.Open(pgDNS), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("error while trying to connect to postgres database")
		return nil, err
	}

	return repo, nil

}

// GetApplications get all valid applications from the database
// deleted applications are ignored
func (pgr *PgRepository) GetApplications(useCache bool) ([]types.Application, error) {
	return nil, nil
}

// GetApplicationsFromCache check if the applications are in cache server
// and return valid items.
func (pgr *PgRepository) GetApplicationsFromCache() ([]types.Application, error) {
	return nil, nil
}

// GetApplication search for a single application with the given ID,
// deleted applications should not return
func (pgr *PgRepository) GetApplication(id uuid.UUID) (types.Application, error) {
	return types.Application{}, nil
}
func (pgr *PgRepository) GetApplicationFromCache(id uuid.UUID) (types.Application, error) {
	return types.Application{}, nil
}

// SaveApplication save a new application in the database
func (pgr *PgRepository) SaveApplication(app types.Application) (uuid.UUID, error) {
	return uuid.New(), nil
}

// UpdateApplication updates attributes from an existing application in the database
// returns an error if its not possible to save data in the database
func (pgr *PgRepository) UpdateApplication(app types.Application) error {
	return nil
}
