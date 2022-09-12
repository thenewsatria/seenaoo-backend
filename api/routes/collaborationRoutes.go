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

	collaborationRoutes.Use(middlewares.CheckAuthorized(userService))
	collaborationRoutes.Post("/", handlers.AddCollaboration(collaborationService, userService, flashcardCoverService))
	collaborationRoutes.Put("/:collaborationId", handlers.UpdateCollaboration(collaborationService))
	collaborationRoutes.Get("/:collaborationId", handlers.AddCollaboration(collaborationService, userService, flashcardCoverService))
	collaborationRoutes.Delete("/:collaborationId", handlers.DeleteCollaboration(collaborationService))
}
