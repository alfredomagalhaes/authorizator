package repository

import (
	"errors"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PgRoleRepository struct {
	DB *gorm.DB
}

// Creates a new role repository
func NewPgRoleRepository(db *gorm.DB) (*PgRoleRepository, error) {
	var roleRepo PgRoleRepository

	if db == nil {
		return nil, errors.New("no database received")
	}

	roleRepo.DB = db
	return &roleRepo, nil
}

// MigrateTable create roles table in the database
func (par PgRoleRepository) MigrateTable() error {
	return par.DB.Migrator().AutoMigrate(types.Role{})
}

// SaveRole save a new application in the database
func (pgr *PgRoleRepository) Save(r types.Role) (uuid.UUID, error) {

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

// Get a specific role record using ID as search key
func (pgr *PgRoleRepository) Get(id uuid.UUID) (types.Role, error) {
	role := types.Role{}

	result := pgr.DB.Where("id = ?", id).First(&role)

	return role, result.Error
}

// Return an Application record using app ID as search key
func (pgr *PgRoleRepository) GetApp(id uuid.UUID) (types.Application, error) {
	appRepo := PgAppRepository{
		DB: pgr.DB,
	}

	return appRepo.Get(id)
}
