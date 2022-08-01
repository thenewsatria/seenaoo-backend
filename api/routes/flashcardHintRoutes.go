package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
)

func flashcardHintRouter(app fiber.Router) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	flashcardHintCollection := database.UseDB().Collection(flashcardhints.CollectionName)
	flashcardRepo := flashcardhints.NewRepo(flashcardHintCollection)
	Service := flashcardhints.NewService(flashcardRepo)

	flashcardHintRoutes.Post("/", handlers.AddFlashcardHint(Service))
}
