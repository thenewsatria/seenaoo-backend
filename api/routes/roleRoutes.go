package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func roleRouter(app fiber.Router, roleService roles.Service, userService users.Service, permissionService permissions.Service) {
	roleRoutes := app.Group("/roles")

	roleRoutes.Use(middlewares.IsLoggedIn(userService))
	roleRoutes.Post("/", handlers.MakeNewRole(roleService))
	roleRoutes.Get("/me", handlers.GetMyRole(roleService))

	//isLoggedIn + only author can see detail, update and delete the role
	roleRoutes.Use(middlewares.IsAuthorized("ROLE", roleService, nil, false, nil))
	roleRoutes.Get("/:roleSlug", handlers.GetRole(roleService, userService, permissionService))
	roleRoutes.Put("/:roleSlug", handlers.UpdateRole(roleService))
	roleRoutes.Delete("/:roleSlug", handlers.DeleteRole(roleService))
}
