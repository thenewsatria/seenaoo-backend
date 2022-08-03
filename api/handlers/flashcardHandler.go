package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddFlashcard(service flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcard := &models.Flashcard{}
		if err := c.BodyParser(flashcard); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_BODY_PARSER_ERROR_MESSAGE))
		}
		result, err := service.InsertFlashcard(flashcard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.FlashcardInsertSuccessResponse(result))
	}
}

func GetFlashcard(flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardId := &models.ReadFlashcardRequest{ID: c.Params("flashcardId")}
		result, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		hints, err := flashcardHintService.PopulateFlashcard(flashcardId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_POPULATE_HINTS_ERROR_MESSAGE))
		}
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardReadSuccessResponse(result, hints))
	}
}
