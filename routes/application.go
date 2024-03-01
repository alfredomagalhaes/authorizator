package routes

import (
	"github.com/alfredomagalhaes/authorizator/controllers"
	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/services"
	"github.com/gofiber/fiber/v2"
)

func ApplicationRoute(route fiber.Router, r repository.AppRepository, ps *services.PermissionService) {

	group := route.Group("application")
	group.Get("/", controllers.GetAllApplications(r))
	group.Get("/:id", controllers.GetApplicationWithID(r))
	group.Get("/:id/roles", controllers.GetAppRoles(ps))
	group.Post("/", controllers.SaveApplication(r))
	group.Put("/:id", controllers.UpdateApplication(r))
}
