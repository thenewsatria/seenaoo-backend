package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func RegisterUser(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hashedUser := c.Locals("hashedUser")
		registeredUser, err := service.InsertUser(hashedUser.(*models.User))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.UserInsertSuccessResponse(registeredUser))
	}
}
