package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func userRouter(app fiber.Router, userProfileService userprofiles.Service, userService users.Service) {
	userRoutes := app.Group("/users")
	// userRoutes.Get("/profile", handlers.)
	userRoutes.Use(middlewares.IsLoggedIn(userService))
	userRoutes.Post("/profile", handlers.EditUserProfile(userProfileService))
}
