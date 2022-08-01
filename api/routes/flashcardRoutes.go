package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
)

func flashcardRouter(app fiber.Router) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))
	flashcardRoutes.Post("/", handlers.AddFlashcard(flashcardService))
}
