package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/refreshtokens"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func authenticationRouter(app fiber.Router, userService users.Service, refreshTokenService refreshtokens.Service) {
	authRoutes := app.Group("/auth")

	authRoutes.Post("/login", handlers.UserLogin(userService, refreshTokenService))
	authRoutes.Post("/register", handlers.RegisterUser(userService, refreshTokenService))
	authRoutes.Post("/token", handlers.RefreshToken(refreshTokenService, userService))

	authRoutes.Use(middlewares.IsLoggedIn(userService))
	authRoutes.Post("/logout", handlers.UserLogout(refreshTokenService))
}
