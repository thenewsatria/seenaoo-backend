package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

func AddFlashcardHint(s flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardHint := &models.FlashcardHint{}
		if err := c.BodyParser(flashcardHint); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.FlashcardHintErrorResponse(err))
		}
		result, err := s.InsertFlashcardHint(flashcardHint)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.FlashcardHintErrorResponse(err))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.FlashcardHintInsertSuccessResponse(result))
	}
}
