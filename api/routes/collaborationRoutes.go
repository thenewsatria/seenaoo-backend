package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func collaborationRouter(app fiber.Router, collaborationService collaborations.Service, userService users.Service,
	userProfileService userprofiles.Service, flashcardCoverService flashcardcovers.Service, roleService roles.Service,
	permissionService permissions.Service) {
	collaborationRoutes := app.Group("/collaboration")

	collaborationRoutes.Use(middlewares.IsLoggedIn(userService))

	//isLoggedIn + only the author can invite the collaboration invites
	collaborationRoutes.Post("/flashcard/:itemId",
		middlewares.IsAllowedToSendCollaboration(flashcardCoverService, collaborationService, roleService, false, true),
		middlewares.HavePermit(permissionService, true, "FLASHCARD.INVITE_COLLABORATOR"),
		handlers.AddCollaboration(collaborationService, userService, flashcardCoverService))

	//isLoggedIn + only author and invited collabotorator can view detail the status (SENT, REJECTED, ACCEPTED)
	collaborationRoutes.Get("/:collaborationId",
		middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, true, true, collaborationService, roleService),
		handlers.GetCollaboration(collaborationService, userService, userProfileService, flashcardCoverService, roleService))

	//isLoggedIn + only invited collabotorator can view detail the status (SENT, REJECTED, ACCEPTED)
	collaborationRoutes.Patch("/:collaborationId",
		middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, true, false, collaborationService, roleService),
		handlers.UpdateCollabStatus(collaborationService))

	//isLoggedIn + only the author can update and delete the collaboration invites=
	collaborationRoutes.Use("/:collaborationId",
		middlewares.IsAuthorized("COLLABORATION", collaborationService, nil, false, true, collaborationService, roleService))

	collaborationRoutes.Delete("/:collaborationId",
		handlers.DeleteCollaboration(collaborationService))

	collaborationRoutes.Put("/:collaborationId",
		handlers.UpdateCollaboration(collaborationService))

}
