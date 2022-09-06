package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func AddCollaboration(collaboratorService collaborations.Service, userService users.Service, flashcardCoverService flashcardCover.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		collaboration := &models.Collaboration{}
		if err := c.BodyParser(collaboration); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_BODY_PARSER_ERROR_MESSAGE))
		}

		collaboration.Inviter = currentUser.Username

		inviterCheckId := &models.UserByUsernameRequest{Username: collaboration.Inviter}
		if !userService.CheckUsernameIsExist(inviterCheckId){
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_INVITER_DOESNT_EXIST_ERROR_MESSAGE))
		}
		
		collaboratorCheckId := &models.UserByUsernameRequest{Username: collaboration.Collaborator}
		if !userService.CheckUsernameIsExist(collaboratorCheckId) {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_COLLABORATOR_DOESNT_EXIST_ERROR_MESSAGE))
		}
		
		switch (collaboration.ItemType) {
		case "FLASHCARD":
			// Check if item is exist
		}
		
		//check if item collaboration is exist

		createdCollaboration, err := service.InsertCollaboration(collaboration)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.CollaborationSuccessResponse(createdCollaboration))
	}
}

func GetCollaboration(collaborationService collaborations.Service, userService users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		collabId := &models.CollaborationById{
			ID: c.Params("collaborationId"),
		}

		collab, err := collaborationService.GetCollaboration(collabId)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_NOT_FOUND_ERROR_MESSAGE))
		}

		userService.CheckUsernameIsExist()
		
		

		c.Status(http.StatusOK)
		return c.JSON(presenters.)
	}
}
