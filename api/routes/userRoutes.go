package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func userRouter(app fiber.Router, userProfileService userprofiles.Service, userService users.Service) {
	userRoutes := app.Group("/my")
	// userRoutes.Get("/profile", handlers.)
	userRoutes.Use(middlewares.IsLoggedIn(userService))
	userRoutes.Get("/profile", handlers.GetUserProfile(userProfileService))
	userRoutes.Put("/profile", handlers.EditUserProfile(userProfileService))
	userRoutes.Delete("/profile/banner", handlers.DeleteProfileBanner(userProfileService))
	userRoutes.Delete("/profile/avatar", handlers.DeleteProfileAvatar(userProfileService))
}
