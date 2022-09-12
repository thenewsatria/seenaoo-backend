package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardRouter(app fiber.Router, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service, userService users.Service) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardRoutes.Use(middlewares.CheckAuthorized(userService))

	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))
	flashcardRoutes.Post("/", handlers.AddFlashcard(flashcardService))
	flashcardRoutes.Put("/:flashcardId", handlers.UpdateFlashcard(flashcardService))
	flashcardRoutes.Delete("/:flashcardId", handlers.DeleteFlashcard(flashcardService))

	flashcardRoutes.Delete("/:flashcardId/delete-all-hints", handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))

}
