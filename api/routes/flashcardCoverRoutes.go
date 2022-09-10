package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardCoverRouter(app fiber.Router, flashcardCoverService flashcardcovers.Service, flashcardService flashcards.Service,
	flashcardHintService flashcardhints.Service, tagService tags.Service, userService users.Service) {

	flashcardCoverRoutes := app.Group("/flashcard-cover")

	//depends on the privacy setting public, unlisted, private
	flashcardCoverRoutes.Get("/:flashcardCoverSlug", handlers.GetFlashcardCover(flashcardCoverService, tagService, userService, flashcardService))

	//isLoggedIn can access
	flashcardCoverRoutes.Post("/", handlers.AddFlashcardCover(flashcardCoverService, tagService))

	//isLoggedIn + author or collaborators can access it
	flashcardCoverRoutes.Put("/:flashcardCoverSlug", handlers.UpdateFlashcardCover(flashcardCoverService, tagService))
	flashcardCoverRoutes.Delete("/:flashcardCoverSlug", handlers.DeleteFlashcardCover(flashcardCoverService))
	flashcardCoverRoutes.Delete("/purge/:flashcardCoverSlug", handlers.PurgeFlashcardCover(flashcardCoverService, flashcardService, flashcardHintService))
}
