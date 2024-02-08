package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ErrNotAllowedChar error = errors.New("use only alphanumeric in role name, to separate strings use /, - ou _")

// CreateRole create a new role in the database, and make some validations
func CreateRole(r repository.Repository, role types.Role) (uuid.UUID, error) {

	var allowedChar = []string{"/", "_", "-"}
	var roleName string
	var regexPattern string = "[A-Za-z0-9]" //just alphanumeric characters

	if role.AppID == uuid.Nil {
		return uuid.Nil, errors.New("role must be associated to an application, 'application_id' missing")
	}

	//Check if application exists
	app, err := r.GetApplication(role.AppID)

	if err != nil {
		return uuid.Nil, errors.New("")
	}

	roleName = role.Name
	//Check role name for forbidden characters
	if strings.Contains(roleName, " ") {
		return uuid.Nil, errors.New("white spaces not allowed in role name")
	}
	for _, char := range allowedChar {
		roleName = strings.ReplaceAll(roleName, char, "")
	}

	match, err := regexp.MatchString(regexPattern, roleName)

	if !match || err != nil {
		return uuid.Nil, ErrNotAllowedChar
	}

	role.Tag = fmt.Sprintf("%s/%s", app.ExternalID, role.Name)

	roleID, err := r.SaveRole(role)

	if err != nil {
		log.Error().Err(err).Msg("failed to create new role in the database")
		return uuid.Nil, errors.New("could not create new role, try again later")
	}

	return roleID, nil
}
