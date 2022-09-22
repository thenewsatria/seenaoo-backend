package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardCoverRouter(app fiber.Router, flashcardCoverService flashcardcovers.Service, flashcardService flashcards.Service,
	flashcardHintService flashcardhints.Service, tagService tags.Service, userService users.Service,
	collaborationService collaborations.Service, roleService roles.Service, permissionService permissions.Service) {

	flashcardCoverRoutes := app.Group("/flashcard-covers")

	flashcardCoverRoutes.Post("/",
		middlewares.IsLoggedIn(userService),
		handlers.AddFlashcardCover(flashcardCoverService, tagService))

	flashcardCoverRoutes.Delete("/purge/:flashcardCoverSlug",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD_COVER", flashcardCoverService, nil, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.PURGE_FLASHCARD"),
		handlers.PurgeFlashcardCover(flashcardCoverService, flashcardService, flashcardHintService))

	//depends on the privacy setting public, unlisted, private
	flashcardCoverRoutes.Get("/:flashcardCoverSlug", handlers.GetFlashcardCover(flashcardCoverService, tagService, userService, flashcardService))

	//isLoggedIn + author or collaborators can access it
	flashcardCoverRoutes.Use("/:flashcardCoverSlug",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD_COVER", flashcardCoverService, nil, true, collaborationService, roleService))

	flashcardCoverRoutes.Put("/:flashcardCoverSlug",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_FLASHCARD"),
		handlers.UpdateFlashcardCover(flashcardCoverService, tagService))

	flashcardCoverRoutes.Delete("/:flashcardCoverSlug",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.DELETE_FLASHCARD"),
		handlers.DeleteFlashcardCover(flashcardCoverService))
}
