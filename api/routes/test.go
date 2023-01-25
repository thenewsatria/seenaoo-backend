package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/tests"
)

func testRouter(app fiber.Router, service tests.Service) {
	testRoutes := app.Group("/tests")
	testRoutes.Post("/", handlers.TestHandler(service))
}
