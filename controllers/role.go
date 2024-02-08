package controllers

import (
	"errors"
	"net/http"

	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/alfredomagalhaes/authorizator/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateRole(r repository.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request types.Role

		err := c.BodyParser(&request)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		roleId, err := usecase.CreateRole(r, request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		locationUrl := buildLocationString(c, roleId.String())
		c.Response().Header.Add("location", locationUrl)
		return c.Status(http.StatusCreated).JSON(SuccessResponse(roleId))
	}

}

func GetRole(r repository.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id", "")

		if id == "" {
			c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("id not informed")))
		}
		parsedID, err := uuid.Parse(id)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("malformed id, check the request")))
		}

		role, err := r.GetRole(parsedID)

		if err != nil {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse(ErrNoRecordsFound))
		}

		return c.Status(http.StatusOK).JSON(SuccessResponse(role))
	}

}
