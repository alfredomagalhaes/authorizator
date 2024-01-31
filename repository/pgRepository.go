package repository

import (
	"fmt"
	"strings"

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
	//postgres://%s:%s@%s:%s/%s?sslmode=disable
	//var pgDNS string = fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s`,
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
	repo.log = log

	return &repo, nil

}

// MigrateTable create all the tables in the database
func (pgr PgRepository) MigrateTables() {
	pgr.DB.Migrator().AutoMigrate(types.Application{})
}

// GetApplications get all valid applications from the database
// deleted applications are ignored
func (pgr *PgRepository) GetApplications(useCache bool) ([]types.Application, error) {

	var appsToRet []types.Application

	result := pgr.DB.Find(&appsToRet)

	return appsToRet, result.Error
}

// GetApplicationsFromCache check if the applications are in cache server
// and return valid items.
func (pgr *PgRepository) GetApplicationsFromCache() ([]types.Application, error) {
	return nil, nil
}

// GetApplication search for a single application with the given ID,
// deleted applications should not return
func (pgr *PgRepository) GetApplication(id uuid.UUID) (types.Application, error) {

	appToRet := types.Application{}

	result := pgr.DB.Where("id = ?", id).First(&appToRet)

	return appToRet, result.Error
}
func (pgr *PgRepository) GetApplicationFromCache(id uuid.UUID) (types.Application, error) {
	return types.Application{}, nil
}

// SaveApplication save a new application in the database
func (pgr *PgRepository) SaveApplication(app types.Application) (uuid.UUID, error) {

	result := pgr.DB.Create(&app)

	if result.Error != nil {
		pgr.log.Error().Err(result.Error).Msg("failed to insert new item in the database")
		if checkIfIsDuplicated(result.Error.Error()) {
			return uuid.Nil, ErrAppDuplicated
		}
		return uuid.Nil, ErrDefaultInsertApp
	}

	return app.ID, nil
}

// UpdateApplication updates attributes from an existing application in the database
// returns an error if its not possible to save data in the database
func (pgr *PgRepository) UpdateApplication(app types.Application) error {

	if app.ID == uuid.Nil {
		return ErrIdToUpdateNill
	}

	result := pgr.DB.Save(&app)

	return result.Error
}

// Check if there is the string "duplicated"
// on error string
func checkIfIsDuplicated(errStr string) bool {
	return strings.Contains(errStr, "duplicate")
}
