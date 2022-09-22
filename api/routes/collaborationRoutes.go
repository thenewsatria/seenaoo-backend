package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func collaborationRouter(app fiber.Router, collaborationService collaborations.Service, userService users.Service,
	flashcardCoverService flashcardcovers.Service, roleService roles.Service, permissionService permissions.Service) {
	collaborationRoutes := app.Group("/collaboration")

	collaborationRoutes.Use(middlewares.IsLoggedIn(userService))

	//isLoggedIn + only the author can invite the collaboration invites
	collaborationRoutes.Post("/flashcard/:itemId",
		middlewares.IsAllowedToSendCollaboration(flashcardCoverService, collaborationService, roleService, false),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.INVITE_COLLABORATOR"),
		handlers.AddCollaboration(collaborationService, userService, flashcardCoverService))

	//isLoggedIn + only author or invited collabotorator can update the status (SENT, REJECTED, ACCEPTED)
	collaborationRoutes.Use("/:collaborationId",
		middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, true, collaborationService, roleService))

	collaborationRoutes.Get("/:collaborationId", handlers.GetCollaboration(collaborationService, userService, flashcardCoverService, roleService))
	collaborationRoutes.Patch("/:collaborationId", handlers.UpdateCollabStatus(collaborationService))

	//isLoggedIn + only the author can update and delete the collaboration invites=
	collaborationRoutes.Use("/:collaborationId",
		middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, false, collaborationService, roleService))

	collaborationRoutes.Delete("/:collaborationId",
		handlers.DeleteCollaboration(collaborationService))

	collaborationRoutes.Put("/:collaborationId",
		handlers.UpdateCollaboration(collaborationService))

}
