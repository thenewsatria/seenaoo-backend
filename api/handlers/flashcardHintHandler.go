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

func AddFlashcardHint(flashcardHintService flashcardhints.Service, flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fcId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		fc, err := flashcardService.FetchFlashcard(fcId)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
		}
		flashcardHint := &models.FlashcardHint{}
		if err := c.BodyParser(flashcardHint); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_BODY_PARSER_ERROR_MESSAGE))
		}
		flashcardHint.FlashcardId = fc.ID
		flashcardHint.FlashcardCoverId = fc.FlashCardCoverId

		result, err, isValidationError := flashcardHintService.InsertFlashcardHint(flashcardHint)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.FlashcardHintSuccessResponse(result))
	}
}

func UpdateFlashcardHint(flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardHintId := &models.FlashcardHintByIdRequest{ID: c.Params("flashcardHintId")}
		flashcardHint, err := flashcardHintService.FetchFlashcardHint(flashcardHintId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		updateBody := &models.FlashcardHint{}
		if err := c.BodyParser(updateBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_BODY_PARSER_ERROR_MESSAGE))
		}

		updateBody.ID = flashcardHint.ID
		updateBody.CreatedAt = flashcardHint.CreatedAt
		updateBody.FlashcardCoverId = flashcardHint.FlashcardCoverId
		updateBody.FlashcardId = flashcardHint.FlashcardId

		updatedFlashcardHint, err, isValidationError := flashcardHintService.UpdateFlashcardHint(updateBody)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON((presenters.ErrorResponse(err.Error())))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardHintSuccessResponse(updatedFlashcardHint))
	}
}

func DeleteFlashcardHint(flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardHintId := &models.FlashcardHintByIdRequest{ID: c.Params("flashcardHintId")}
		flashcardHint, err := flashcardHintService.FetchFlashcardHint(flashcardHintId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		deletedFlashcardHint, err := flashcardHintService.RemoveFlashcardHint(flashcardHint)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardHintSuccessResponse(deletedFlashcardHint))
	}
}

func DeleteFlashcardHintsByFlashcardId(flashcardHintService flashcardhints.Service, flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		_, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		deletedFHintsCount, err := flashcardHintService.RemoveFlashcardHintsByFlashcardId(flashcardId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		if deletedFHintsCount == 0 {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HAS_EMPTY_HINTS_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.BatchOperationResponse("DELETE", "flashcard hints", deletedFHintsCount))
	}
}
