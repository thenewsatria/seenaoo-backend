package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddFlashcard(service flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcard := &models.Flashcard{}
		if err := c.BodyParser(flashcard); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.FlashcardErrorResponse(err))
		}
		result, err := service.InsertFlashcard(flashcard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.FlashcardErrorResponse(err))
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
				return c.JSON(presenters.FlashCardNotFound())
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.FlashcardErrorResponse(err))
		}

		hints, err := flashcardHintService.PopulateFlashcard(flashcardId)
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardReadSuccessResponse(result, hints))
	}
}
