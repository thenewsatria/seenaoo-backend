package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		fcCoverRequest := &models.FlashcardCoverRequest{}
		err := c.BodyParser(fcCoverRequest)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		tagIds := []primitive.ObjectID{}

		for _, tagString := range fcCoverRequest.Tags {
			tagName := &models.TagByName{TagName: tagString}
			existedTag, err := tagService.FetchTagByName(tagName)
			if err != nil {
				if err == mongo.ErrNoDocuments { //jika tag tidak ada maka buat baru
					tag := &models.Tag{TagName: tagString}
					newTag, err := tagService.InsertTag(tag)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_INSERT_ERROR_MESSAGE))
					}
					tagIds = append(tagIds, newTag.ID)
					continue
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			tagIds = append(tagIds, existedTag.ID)
		}

		currentTimeStr := fmt.Sprintf("%v", time.Now().Unix())

		fcCover := &models.FlashcardCover{
			Slug:        slug.Make(fcCoverRequest.Title) + "-" + currentTimeStr,
			Title:       fcCoverRequest.Title,
			Description: fcCoverRequest.Description,
			Image_path:  fcCoverRequest.Image_path,
			Tags:        tagIds,
			Author:      currentUser.Username,
		}

		insertedFcCover, err := flashcardCoverService.InsertFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		return c.JSON(presenters.FlashcardCoverSuccessResponse(insertedFcCover))
	}
}

func GetFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service, userService users.Service, flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		tagDetails := []models.Tag{}

		for _, fcTag := range fcCover.Tags {
			tagId := &models.TagById{ID: fcTag.Hex()}
			tag, err := tagService.FetchTagById(tagId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.TAG_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			tagDetails = append(tagDetails, *tag)
		}

		userUname := &models.UserByUsernameRequest{Username: fcCover.Author}
		author, err := userService.FetchUserByUsername(userUname)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		fcCvrId := &models.FlashcardCoverById{ID: fcCover.ID.Hex()}
		flashcards, err := flashcardService.PopulateFlashcardCover(fcCvrId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_POPULATE_FLASHCARDS_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverDetailSuccessResponse(fcCover, &tagDetails, flashcards, author))
	}
}

func UpdateFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		updateBody := &models.FlashcardCoverRequest{}
		if err := c.BodyParser(updateBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		newSlug := slug.Make(updateBody.Title) + "-" + fmt.Sprintf("%v", time.Now().Unix())
		fcCover.Slug = newSlug
		fcCover.Title = updateBody.Title
		fcCover.Description = updateBody.Description
		fcCover.Image_path = updateBody.Image_path

		tagIds := []primitive.ObjectID{}

		for _, tagString := range updateBody.Tags {
			tagName := &models.TagByName{TagName: tagString}
			existedTag, err := tagService.FetchTagByName(tagName)
			if err != nil {
				if err == mongo.ErrNoDocuments { //jika tag tidak ada maka buat baru
					tag := &models.Tag{TagName: tagString}
					newTag, err := tagService.InsertTag(tag)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_INSERT_ERROR_MESSAGE))
					}
					tagIds = append(tagIds, newTag.ID)
					continue
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			tagIds = append(tagIds, existedTag.ID)
		}

		fcCover.Tags = tagIds
		updatedFcCover, err := flashcardCoverService.UpdateFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(updatedFcCover))
	}
}

func DeleteFlashcardCover(flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
		}

		deletedFcCover, err := flashcardCoverService.RemoveFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(deletedFcCover))
	}
}

func PurgeFlashcardCover(flashcardCoverService flashcardcovers.Service, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//Get Flashcard Cover
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
		}

		//Get flashcards by flashcard cover id so we able to delete each flashcard's hints
		fcCoverId := &models.FlashcardCoverById{ID: fcCover.ID.Hex()}
		flashcards, err := flashcardService.PopulateFlashcardCover(fcCoverId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_POPULATE_FLASHCARDS_ERROR_MESSAGE))
		}

		//Loop each flashcard to delete each flashcard's hints
		for _, flashcard := range *flashcards {
			flashcardId := &models.FlashcardByIdRequest{ID: flashcard.ID.Hex()}
			_, err := flashcardHintService.RemoveFlashcardHintsByFlashcardId(flashcardId)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE))
			}
		}

		//Delete all flashcard with the same flashcard cover id
		_, err = flashcardService.RemoveFlashcardsByFlashcardCoverId(fcCoverId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// Delete The flashcard cover
		deletedFcCover, err := flashcardCoverService.RemoveFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(deletedFcCover))
	}
}
