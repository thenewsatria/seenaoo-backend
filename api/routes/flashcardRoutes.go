package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
)

func flashcardRouter(app fiber.Router, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))
	flashcardRoutes.Post("/", handlers.AddFlashcard(flashcardService))
}
