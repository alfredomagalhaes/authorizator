package services

import (
	"errors"

	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
)

type PermissionService struct {
	appRepo  repository.AppRepository
	roleRepo repository.RoleRepository
}

var ErrAppRepositoryNil error = errors.New("application repository not received")
var ErrRoleRepositoryNil error = errors.New("role repository not received")

// Initialize a permission service struct
// to let use different method from the repositories
func NewPermissionService(appRepo repository.AppRepository, roleRepo repository.RoleRepository) (*PermissionService, error) {
	var permService PermissionService

	if appRepo == nil {
		return &permService, ErrAppRepositoryNil
	}

	permService.appRepo = appRepo

	if roleRepo == nil {
		return &permService, ErrRoleRepositoryNil
	}

	permService.roleRepo = roleRepo

	return &permService, nil

}

func (ps PermissionService) GetAppWithRoles(appID uuid.UUID) (types.Application, error) {
	//Check if app exists
	app, err := ps.GetApp(appID)

	//Get roles associated with app
	roles, err := ps.roleRepo.GetAppRoles(app.ID)

	app.Roles = roles
	return app, err
}

func (ps PermissionService) GetApp(appID uuid.UUID) (types.Application, error) {
	//Check if app exists
	app, err := ps.appRepo.Get(appID)

	if err != nil {

	}

	return app, nil
}
