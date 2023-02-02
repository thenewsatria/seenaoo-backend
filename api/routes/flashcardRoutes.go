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
	flashcardRoutes := app.Group("/flashcards")

	flashcardRoutes.Post("/add/:flashcardCoverSlug",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD_COVER", flashcardCoverService, nil, true, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.ADD_CARD"),
		handlers.AddFlashcard(flashcardService, flashcardCoverService))

	flashcardRoutes.Delete("/purge/:flashcardId",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.PURGE_CARD"),
		handlers.PurgeFlashcard(flashcardService, flashcardHintService))

	flashcardRoutes.Delete("/clear-hints/:flashcardId",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.CLEAR_CARD_HINT"),
		handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))

	flashcardRoutes.Delete("/frontimage/:flashcardId",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_CARD"),
		handlers.DeleteFlashcardFrontImage(flashcardService))

	flashcardRoutes.Delete("/backimage/:flashcardId",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_CARD"),
		handlers.DeleteFlashcardBackImage(flashcardService))

	//depends on the cover privacy setting
	flashcardRoutes.Get("/:flashcardId", handlers.GetFlashcard(flashcardService, flashcardHintService))

	flashcardRoutes.Use("/:flashcardId",
		middlewares.IsLoggedIn(userService),
		middlewares.IsAuthorized("FLASHCARD", flashcardService, flashcardCoverService, true, true, collaborationService, roleService))

	flashcardRoutes.Put("/:flashcardId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_CARD"),
		handlers.UpdateFlashcard(flashcardService))

	flashcardRoutes.Delete("/:flashcardId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.DELETE_CARD"),
		handlers.DeleteFlashcard(flashcardService))

}
