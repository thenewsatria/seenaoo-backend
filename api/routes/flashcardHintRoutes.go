package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
)

func flashcardHintRouter(app fiber.Router, flashcardHintService flashcardhints.Service, flashcardService flashcards.Service) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	flashcardHintRoutes.Post("/", handlers.AddFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Put("/:flashcardHintId", handlers.UpdateFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Delete("/:flashcardHintId", handlers.DeleteFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Delete("/deletes-by-flashcard-id/:flashcardId", handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))
}
