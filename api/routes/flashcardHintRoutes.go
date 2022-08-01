package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
)

func flashcardHintRouter(app fiber.Router, flashcardHintService flashcardhints.Service) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	flashcardHintRoutes.Post("/", handlers.AddFlashcardHint(flashcardHintService))
}
