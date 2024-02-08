package routes

import (
	"github.com/alfredomagalhaes/authorizator/controllers"
	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/gofiber/fiber/v2"
)

func RolesRoute(route fiber.Router, r repository.Repository) {

	group := route.Group("role")
	group.Post("/", controllers.CreateRole(r))
	group.Get("/:id", controllers.GetRole(r))

}
