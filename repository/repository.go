package repository

import (
	"errors"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
)

var ErrNoRecordsFound = errors.New("no records found with given parameters")
var ErrAppDuplicated = errors.New("application already exists, try another `external_id`")
var ErrRoleDuplicated = errors.New("role already exists, try another `name`")
var ErrDefaultInsertApp = errors.New("error while trying to create an application, try again later")
var ErrDefaultInsertRole = errors.New("error while trying to create a role, try again later")
var ErrIdToUpdateNill = errors.New("id to update can't be nill")

type Repository interface {
	GetApplications(useCache bool) ([]types.Application, error)
	GetApplicationsFromCache() ([]types.Application, error)
	GetApplication(id uuid.UUID) (types.Application, error)
	GetApplicationFromCache(id uuid.UUID) (types.Application, error)
	SaveApplication(app types.Application) (uuid.UUID, error)
	UpdateApplication(app types.Application) error
	SaveRole(r types.Role) (uuid.UUID, error)
	GetRole(uuid.UUID) (types.Role, error)
}
