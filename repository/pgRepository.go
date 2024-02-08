package repository

import (
	"fmt"
	"strings"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
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

// MigrateTable create all the tables in the database
func (pgr PgRepository) MigrateTables() {
	//pgr.DB.Migrator().AutoMigrate(types.Application{})
	//pgr.DB.Migrator().AutoMigrate(types.Role{})
	//pgr.DB.AutoMigrate(types.Application{})
	pgr.DB.Migrator().AutoMigrate(types.Application{})
	pgr.DB.Migrator().AutoMigrate(types.Role{})

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
		log.Error().Err(result.Error).Msg("failed to insert new item in the database")
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

// SaveApplication save a new application in the database
func (pgr *PgRepository) SaveRole(r types.Role) (uuid.UUID, error) {

	result := pgr.DB.Create(&r)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to insert new item in the database")
		if checkIfIsDuplicated(result.Error.Error()) {
			return uuid.Nil, ErrRoleDuplicated
		}
		return uuid.Nil, ErrDefaultInsertRole
	}

	return r.ID, nil
}

func (pgr *PgRepository) GetRole(id uuid.UUID) (types.Role, error) {
	role := types.Role{}

	result := pgr.DB.Where("id = ?", id).First(&role)

	return role, result.Error
}
