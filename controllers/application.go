package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var ErrNoRecordsFound error = errors.New("no records found")
var ErrMalformedRequest error = errors.New("malformed body, check the request")
var ErrIDNotInformed error = errors.New("id not informed")
var ErrMalformedID error = errors.New("malformed id, check the request")

// GetAllApplications returns all the applications
// stored in the database
func GetAllApplications(r repository.AppRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		appList, err := r.GetAll(true)

		if err != nil {
			c.Status(http.StatusNotFound).JSON(ErrorResponse(err))
		}

		if len(appList) == 0 {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse(ErrNoRecordsFound))
		}

		return c.Status(http.StatusOK).JSON(SuccessResponse(appList))
	}
}

// GetApplicationWithID search for an application
// based on the id from the request
func GetApplicationWithID(r repository.AppRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id, err := GetIDFromRequest(c)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(err))
		}
		app, err := r.Get(id)

		if err != nil {
			return c.Status(http.StatusNotFound).JSON(ErrorResponse(ErrNoRecordsFound))
		}

		return c.Status(http.StatusOK).JSON(SuccessResponse(app))

	}
}

// SaveApplication creates a new application in database
func SaveApplication(r repository.AppRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request types.Application

		err := c.BodyParser(&request)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		insertedID, err := r.Save(request)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err == repository.ErrAppDuplicated {
				statusCode = http.StatusBadRequest
			}
			return c.Status(statusCode).JSON(ErrorResponse(err))
		}
		locationUrl := buildLocationString(c, insertedID.String())
		c.Response().Header.Add("location", locationUrl)
		return c.Status(http.StatusCreated).JSON(SuccessResponse(insertedID))
	}
}

func UpdateApplication(r repository.AppRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request types.Application

		id, err := GetIDFromRequest(c)

		if err != nil {
			c.Status(http.StatusBadRequest).JSON(ErrorResponse(err))
		}

		err = c.BodyParser(&request)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrMalformedRequest))
		}

		request.ID = id

		err = r.Update(request)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("failed to perform update")))
		}

		return c.Status(http.StatusOK).JSON(SuccessResponse("ok"))
	}
}

// ErrorResponse creates a default error response
// struct, to be used in all api's
func ErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"success": false,
		"error":   err.Error(),
	}
}

// SuccessResponse creates a default success response
// struct, to be used in all api's
func SuccessResponse(v any) *fiber.Map {
	return &fiber.Map{
		"success": true,
		"data":    v,
	}
}

func buildLocationString(c *fiber.Ctx, id string) string {
	var path string = fmt.Sprintf("%s/%s", string(c.Request().URI().Path()), id)

	return path
}

// Check id from url params and try to convert
// to an uuid
func GetIDFromRequest(c *fiber.Ctx) (uuid.UUID, error) {
	var nilID uuid.UUID = uuid.Nil
	var parseID uuid.UUID
	id := c.Params("id")

	if id == "" {
		return nilID, ErrIDNotInformed
	}

	parseID, err := uuid.Parse(id)

	if err != nil {
		return nilID, ErrMalformedID
	}

	return parseID, nil
}
