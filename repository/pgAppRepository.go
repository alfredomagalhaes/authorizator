package repository

import (
	"errors"
	"strings"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Struct to contain database connections
// to be used in applications repository interface
// using postgres as database server
type PgAppRepository struct {
	DB *gorm.DB
}

func NewPgAppRepository(db *gorm.DB) (*PgAppRepository, error) {
	var appRepo PgAppRepository

	if db == nil {
		return nil, errors.New("no database received")
	}

	appRepo.DB = db
	return &appRepo, nil
}

// MigrateTable create applications table in the database
func (par PgAppRepository) MigrateTable() error {
	return par.DB.Migrator().AutoMigrate(types.Application{})
}

// GetAll get all valid applications from the database
// deleted applications are ignored
func (par PgAppRepository) GetAll(useCache bool) ([]types.Application, error) {

	var appsToRet []types.Application

	result := par.DB.Find(&appsToRet)

	return appsToRet, result.Error
}

// GetAllFromCache check if the applications are in cache server
// and return valid items.
func (par PgAppRepository) GetAllFromCache() ([]types.Application, error) {
	return nil, nil
}

// Get search for a single application with the given ID,
// deleted applications should not return
func (par PgAppRepository) Get(id uuid.UUID) (types.Application, error) {

	appToRet := types.Application{}

	result := par.DB.Where("id = ?", id).First(&appToRet)

	return appToRet, result.Error
}
func (par PgAppRepository) GetFromCache(id uuid.UUID) (types.Application, error) {
	return types.Application{}, nil
}

// Save save a new application in the database
func (par PgAppRepository) Save(app types.Application) (uuid.UUID, error) {

	result := par.DB.Create(&app)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to insert new item in the database")
		if checkIfIsDuplicated(result.Error.Error()) {
			return uuid.Nil, ErrAppDuplicated
		}
		return uuid.Nil, ErrDefaultInsertApp
	}

	return app.ID, nil
}

// Update updates attributes from an existing application in the database
// returns an error if its not possible to save data in the database
func (par PgAppRepository) Update(app types.Application) error {

	if app.ID == uuid.Nil {
		return ErrIdToUpdateNill
	}

	result := par.DB.Save(&app)

	return result.Error
}

// Check if there is the string "duplicated"
// on error string
func checkIfIsDuplicated(errStr string) bool {
	return strings.Contains(errStr, "duplicate")
}
