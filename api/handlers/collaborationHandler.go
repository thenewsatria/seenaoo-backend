package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddCollaboration(collaboratorService collaborations.Service, userService users.Service, flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		collaboration := &models.Collaboration{}
		if err := c.BodyParser(collaboration); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_BODY_PARSER_ERROR_MESSAGE))
		}

		collaboration.Inviter = currentUser.Username

		inviterCheckId := &models.UserByUsernameRequest{Username: collaboration.Inviter}
		if !userService.CheckUsernameIsExist(inviterCheckId) {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_INVITER_DOESNT_EXIST_ERROR_MESSAGE))
		}

		collaboratorCheckId := &models.UserByUsernameRequest{Username: collaboration.Collaborator}
		if !userService.CheckUsernameIsExist(collaboratorCheckId) {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_COLLABORATOR_DOESNT_EXIST_ERROR_MESSAGE))
		}

		switch collaboration.ItemType {
		case "FLASHCARD":
			fcCoverId := &models.FlashcardCoverById{ID: collaboration.ItemID.Hex()}
			_, err := flashcardCoverService.FetchFlashcardCoverById(fcCoverId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
		default:
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ITEM_TYPE_IS_UNKNOWN))
		}

		createdCollaboration, err := collaboratorService.InsertCollaboration(collaboration)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.CollaborationSuccessResponse(createdCollaboration))
	}
}

func GetCollaboration(collaborationService collaborations.Service, userService users.Service, flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		collabId := &models.CollaborationById{
			ID: c.Params("collaborationId"),
		}
		collab, err := collaborationService.GetCollaboration(collabId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		invUsername := &models.UserByUsernameRequest{Username: collab.Inviter}
		inviter, err := userService.FetchUserByUsername(invUsername)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		cbtrUsername := &models.UserByUsernameRequest{Username: collab.Collaborator}
		collaborator, err := userService.FetchUserByUsername(cbtrUsername)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		switch collab.ItemType {
		case "FLASHCARD":
			fcId := &models.FlashcardCoverById{ID: collab.ItemID.Hex()}
			flashcardCvr, err := flashcardCoverService.FetchFlashcardCoverById(fcId)
			if err != nil {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}

			c.Status(http.StatusOK)
			return c.JSON(presenters.CollaborationFlashcardDetailSuccessResponse(collab, inviter, collaborator, flashcardCvr))
		default:
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ITEM_TYPE_IS_UNKNOWN))
		}
	}
}

func UpdateCollaboration(service collaborations.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		collabId := &models.CollaborationById{ID: c.Params("collaborationId")}
		collab, err := service.GetCollaboration(collabId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		updateBody := &models.Collaboration{}
		err = c.BodyParser(updateBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_BODY_PARSER_ERROR_MESSAGE))
		}

		updateBody.ID = collab.ID
		updatedCollab, err := service.UpdateCollaboration(updateBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.CollaborationSuccessResponse(updatedCollab))
	}
}

func DeleteCollaboration(service collaborations.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		collabId := &models.CollaborationById{ID: c.Params("collaborationId")}
		collab, err := service.GetCollaboration(collabId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		deleted, err := service.RemoveCollaboration(collab)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.CollaborationSuccessResponse(deleted))
	}
}
