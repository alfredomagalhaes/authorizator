package repository

import (
	"errors"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
)

var ErrNoRecordsFound = errors.New("no records found with given parameters")
var ErrAppDuplicated = errors.New("application already exists, try another `external_id")
var ErrDefaultInsertApp = errors.New("error while trying to create an application, try again later")

type Repository interface {
	GetApplications(useCache bool) ([]types.Application, error)
	GetApplicationsFromCache() ([]types.Application, error)
	GetApplication(id uuid.UUID) (types.Application, error)
	GetApplicationFromCache(id uuid.UUID) (types.Application, error)
	SaveApplication(app types.Application) (uuid.UUID, error)
	UpdateApplication(app types.Application) error
}
