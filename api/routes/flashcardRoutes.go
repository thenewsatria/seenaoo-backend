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
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func flashcardRouter(app fiber.Router, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service,
	flashcardCoverService flashcardcovers.Service, userService users.Service, collaborationService collaborations.Service,
	roleService roles.Service, permissionService permissions.Service) {
	flashcardRoutes := app.Group("/flashcard")

	//depends on the cover privacy setting
	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))

	//isLoggedIn + author or collaborators can access
	flashcardRoutes.Use(middlewares.IsLoggedIn(userService))

	flashcardRoutes.Post("/add/:flashcardCoverSlug",
		middlewares.IsAuthorized("FLASHCARD_COVER", flashcardCoverService, nil, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.ADD_CARD"),
		handlers.AddFlashcard(flashcardService, flashcardCoverService))

	flashcardRoutes.Use("/:flashcardId",
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, collaborationService, roleService))

	flashcardRoutes.Put("/:flashcardId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_CARD"),
		handlers.UpdateFlashcard(flashcardService))

	flashcardRoutes.Delete("/:flashcardId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.DELETE_CARD"),
		handlers.DeleteFlashcard(flashcardService))

	flashcardRoutes.Delete("/purge/:flashcardId",
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.PURGE_CARD"),
		handlers.PurgeFlashcard(flashcardService, flashcardHintService))

	flashcardRoutes.Delete("/:flashcardId/delete-all-hints",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.CLEAR_CARD_HINT"),
		handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))
}
