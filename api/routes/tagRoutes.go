package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/handlers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
)

func tagRouter(app fiber.Router, tagService tags.Service, flashcardCoverService flashcardcovers.Service) {
	tagRoutes := app.Group("/tags")

	//accessible for everyone
	tagRoutes.Get("/:tagName", handlers.ShowItemsByTagName(tagService, flashcardCoverService))
}
