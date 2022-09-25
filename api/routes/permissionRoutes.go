package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
)

func permissionRouter(app fiber.Router, permissionService permissions.Service) {
	permissionRoutes := app.Group("/permissions")

	permissionRoutes.Get("/categories", handlers.GetAvailablePermissions(permissionService))
}
