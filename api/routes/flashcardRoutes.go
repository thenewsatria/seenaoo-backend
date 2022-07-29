package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
)

func flashcardRouter(app fiber.Router) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardRoutes.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "success",
		})
	})

	flashcardRoutes.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "success",
		})
	})

	flashcardRoutes.Post("/", handlers.CreateFlashcard)
}
