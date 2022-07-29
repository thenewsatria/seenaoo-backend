package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
)

func flashcardRouter(app fiber.Router) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardCollection := database.UseDB().Collection(flashcards.CollectionName)
	flashcardRepo := flashcards.NewRepo(flashcardCollection)
	service := flashcards.NewService(flashcardRepo)

	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(service))
	flashcardRoutes.Post("/", handlers.AddBook(service))
}
