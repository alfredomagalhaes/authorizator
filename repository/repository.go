package repository

import (
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
)

type Repository interface {
	GetApplications(useCache bool) ([]types.Application, error)
	GetApplicationsFromCache() ([]types.Application, error)
	GetApplication(id uuid.UUID) (types.Application, error)
	GetApplicationFromCache(id uuid.UUID) (types.Application, error)
	SaveApplication(app types.Application) (uuid.UUID, error)
}
