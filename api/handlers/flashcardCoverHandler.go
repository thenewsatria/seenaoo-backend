package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
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

		_, err = flashcardCoverService.InsertFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		return c.JSON(presenters.FlashcardCoverSuccessResponse(fcCover))
	}
}

func GetFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service, userService users.Service, flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fcCoverId := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverId)
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
