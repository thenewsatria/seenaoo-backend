package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
)

func flashcardRouter(app fiber.Router) {
	flashcardRoutes := app.Group("/flashcard")

	flashcardCollection := database.UseDB().Collection(flashcards.CollectionName)
	flashcardRepo := flashcards.NewRepo(flashcardCollection)
	flashcardService := flashcards.NewService(flashcardRepo)

	flashcardHintCollection := database.UseDB().Collection(flashcardhints.CollectionName)
	flashcardHintRepo := flashcardhints.NewRepo(flashcardHintCollection)
	flashcardHintService := flashcardhints.NewService(flashcardHintRepo)

	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))
	flashcardRoutes.Post("/", handlers.AddFlashcard(flashcardService))
}
