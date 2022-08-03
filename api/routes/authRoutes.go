package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func authenticationRouter(app fiber.Router, userService users.Service) {
	authRoutes := app.Group("/auth")

	hashPasswordMiddlewareRoutes := authRoutes.Use(middlewares.HashUserPassword())
	hashPasswordMiddlewareRoutes.Post("/register", handlers.RegisterUser(userService))

	authRoutes.Post("/login", handlers.UserLogin(userService))
}
