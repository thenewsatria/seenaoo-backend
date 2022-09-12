package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardHintRouter(app fiber.Router, flashcardHintService flashcardhints.Service, flashcardService flashcards.Service, userService users.Service) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	flashcardHintRoutes.Use(middlewares.CheckAuthorized(userService))

	flashcardHintRoutes.Post("/", handlers.AddFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Put("/:flashcardHintId", handlers.UpdateFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Delete("/:flashcardHintId", handlers.DeleteFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Delete("/deletes-by-flashcard-id/:flashcardId", handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))
}
