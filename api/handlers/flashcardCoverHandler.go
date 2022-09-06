package handlers

import (
	"fmt"
	"net/http"
	"slug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func AddFlashcardCover(service flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		fcCoverRequest := &models.FlashcardCoverRequest{}
		err := c.BodyParser(fcCoverRequest)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}
		
		currentTimeStr := fmt.Sprintf("%v", time.Now().Unix())
		fcCover.Slug = 
		fcCover := &models.FlashcardCover{
			Slug: slug.Make(fcCover.Title) + "-" + currentTimeStr,
		}
		fcCover.Author = currentUser.ID
		

		_, err = service.InsertFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		return c.JSON(presenters.FlashcardCoverSuccessResponse(fcCover))
	}
}
