package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func userRouter(app fiber.Router, userService users.Service) {
	userRoutes := app.Group("/user")

	test := userRoutes.Use(middlewares.HashUserPassword())
	test.Post("/register", handlers.RegisterUser(userService))
}
