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

func flashcardHintRouter(app fiber.Router, flashcardHintService flashcardhints.Service, flashcardService flashcards.Service,
	flashcardCoverService flashcardcovers.Service, userService users.Service, collaborationService collaborations.Service,
	roleService roles.Service, permissionService permissions.Service) {
	flashcardHintRoutes := app.Group("/flashcard-hint")

	//isLoggedIn + author or collaborators can access
	flashcardHintRoutes.Use(middlewares.IsLoggedIn(userService))

	//need to FIX this add flashcard hint authorization
	flashcardHintRoutes.Post("/",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.ADD_CARD_HINT"),
		handlers.AddFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Use("/:flashcardHintId",
		middlewares.IsAuthorized("FLASHCARD_HINT", flashcardHintService, flashcardCoverService, true, collaborationService, roleService))

	flashcardHintRoutes.Put("/:flashcardHintId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.UPDATE_CARD_HINT"),
		handlers.UpdateFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Delete("/:flashcardHintId",
		middlewares.HavePermit(permissionService, true, "FLASHCARD.DELETE_CARD_HINT"),
		handlers.DeleteFlashcardHint(flashcardHintService))

	flashcardHintRoutes.Delete("/deletes-by-flashcard-id/:flashcardId",
		middlewares.IsAuthorized("FLASHCARD", flashcardHintService, flashcardCoverService, true, collaborationService, roleService),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.CLEAR_CARD_PERMIT"),
		handlers.DeleteFlashcardHintsByFlashcardId(flashcardHintService, flashcardService))
}
