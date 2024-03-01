package controllers

import (
	"net/http"

	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/services"
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/gofiber/fiber/v2"
)

// Controller to create roles from the rest calls
func CreateRole(r repository.RoleRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request types.Role

		err := c.BodyParser(&request)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		roleId, err := services.CreateRole(r, request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		locationUrl := buildLocationString(c, roleId.String())
		c.Response().Header.Add("location", locationUrl)
		return c.Status(http.StatusCreated).JSON(SuccessResponse(roleId))
	}

}

// Controller to return a existing role filtered by ID
func GetRole(r repository.RoleRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := GetIDFromRequest(c)

		if err != nil {
			c.Status(http.StatusBadRequest).JSON(ErrorResponse(err))
		}

		role, err := r.Get(id)

		if err != nil {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse(ErrNoRecordsFound))
		}

		return c.Status(http.StatusOK).JSON(SuccessResponse(role))
	}

}
