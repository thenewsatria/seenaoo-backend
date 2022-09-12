package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func collaborationRouter(app fiber.Router, collaborationService collaborations.Service, userService users.Service, flashcardCoverService flashcardcovers.Service) {
	collaborationRoutes := app.Group("/collaboration")

	collaborationRoutes.Use(middlewares.IsLoggedIn(userService))

	//isLoggedIn + only the author can invite the collaboration invites
	collaborationRoutes.Post("/flashcard/:itemId",
		middlewares.IsAllowedToSendCollaboration(flashcardCoverService, collaborationService, false),
		handlers.AddCollaboration(collaborationService, userService, flashcardCoverService))

	//isLoggedIn + only the author can delete the collaboration invites
	collaborationRoutes.Delete("/:collaborationId", middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, false, collaborationService),
		handlers.DeleteCollaboration(collaborationService))

	//isLoggedIn + only author or invited collabotorator can update the status
	collaborationRoutes.Use(middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, true, collaborationService))
	collaborationRoutes.Put("/:collaborationId", handlers.UpdateCollaboration(collaborationService))
	collaborationRoutes.Get("/:collaborationId", handlers.GetCollaboration(collaborationService, userService, flashcardCoverService))

}
