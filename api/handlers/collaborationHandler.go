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

func AddCollaboration(collaboratorService collaborations.Service, userService users.Service, itemService interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//get current user
		currentUser := c.Locals("currentUser").(*models.User)

		//Parse collaboration body request
		collaboration := &models.Collaboration{}
		if err := c.BodyParser(collaboration); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_BODY_PARSER_ERROR_MESSAGE))
		}

		//Set inviter as current logged in user username
		collaboration.Inviter = currentUser.Username

		inviterCheckUname := &models.UserByUsernameRequest{Username: collaboration.Inviter}
		if !userService.CheckUsernameIsExist(inviterCheckUname) {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_INVITER_DOESNT_EXIST_ERROR_MESSAGE))
		}

		collaboratorCheckUname := &models.UserByUsernameRequest{Username: collaboration.Collaborator}
		if !userService.CheckUsernameIsExist(collaboratorCheckUname) {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_COLLABORATOR_DOESNT_EXIST_ERROR_MESSAGE))
		}

		//Checking ownership of item will be done in middleware
		switch collaboration.ItemType {
		case "FLASHCARD":
			fcCoverId := &models.FlashcardCoverById{ID: c.Params("itemId")}
			flashcardCoverService := itemService.(flashcardcovers.Service)
			fcCvr, err := flashcardCoverService.FetchFlashcardCoverById(fcCoverId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			collaboration.ItemID = fcCvr.ID
		default:
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ITEM_TYPE_IS_UNKNOWN))
		}

		//if previous collaboration already exist then edit the entity instead, otherwise create new collaboration
		cItemIdAndCollaborator := &models.CollaborationByItemIdAndCollaborator{
			ItemID:       collaboration.ItemID.Hex(),
			Collaborator: collaboration.Collaborator,
		}
		existedCollab, err := collaboratorService.FetchCollaborationByItemIdAndCollaborator(cItemIdAndCollaborator)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				createdCollaboration, err := collaboratorService.InsertCollaboration(collaboration)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_INSERT_ERROR_MESSAGE))
				}
				c.Status(http.StatusOK)
				return c.JSON(presenters.CollaborationSuccessResponse(createdCollaboration))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		if existedCollab.Status == "REJECTED" {
			existedCollab.Status = "SENT"
			updatedCollab, err := collaboratorService.UpdateCollaboration(existedCollab)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_UPDATE_ERROR_MESSAGE))
			}
			c.Status(http.StatusOK)
			return c.JSON(presenters.CollaborationSuccessResponse(updatedCollab))
		} else {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ALREADY_EXIST_ERROR_MESSAGE))
		}
	}
}

func GetCollaboration(collaborationService collaborations.Service, userService users.Service, flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		collabId := &models.CollaborationById{
			ID: c.Params("collaborationId"),
		}
		collab, err := collaborationService.FetchCollaboration(collabId)
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
		collab, err := service.FetchCollaboration(collabId)
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
		collab, err := service.FetchCollaboration(collabId)
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
