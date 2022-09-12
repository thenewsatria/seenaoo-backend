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

func flashcardRouter(app fiber.Router, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service,
	flashcardCoverService flashcardcovers.Service, userService users.Service, collaborationService collaborations.Service) {
	flashcardRoutes := app.Group("/flashcard")

	//depends on the cover privacy setting
	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))

	//isLoggedIn + author or collaborators can access
	flashcardRoutes.Use(middlewares.IsLoggedIn(userService))
	flashcardRoutes.Use(middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, collaborationService))

	flashcardRoutes.Post("/", handlers.AddFlashcard(flashcardService))
	flashcardRoutes.Put("/:flashcardId", handlers.UpdateFlashcard(flashcardService))
	flashcardRoutes.Delete("/:flashcardId", handlers.DeleteFlashcard(flashcardService))

	flashcardRoutes.Delete("/:flashcardId/delete-all-hints", handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))

}
