package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "success",
		})
	})

	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	flashcardRouter(apiV1)
	flashcardHintRouter(apiV1)
}
