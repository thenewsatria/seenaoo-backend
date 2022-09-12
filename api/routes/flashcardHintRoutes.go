package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardHintRouter(app fiber.Router, flashcardHintService flashcardhints.Service, flashcardService flashcards.Service,
	flashcardCoverService flashcardcovers.Service, userService users.Service, collaborationService collaborations.Service) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	//isLoggedIn + author or collaborators can access
	flashcardHintRoutes.Use(middlewares.IsLoggedIn(userService))
	flashcardHintRoutes.Use(middlewares.IsAuthorized("FLASHCARD_HINT", flashcardHintService, flashcardCoverService, true, collaborationService))

	flashcardHintRoutes.Post("/", handlers.AddFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Put("/:flashcardHintId", handlers.UpdateFlashcardHint(flashcardHintService))
	flashcardHintRoutes.Delete("/:flashcardHintId", handlers.DeleteFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Delete("/deletes-by-flashcard-id/:flashcardId", handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))
}
