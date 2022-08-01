package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func RegisterUser(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hashedUser := c.Locals("hashedUser").(*models.User)
		userEmail := &models.UserByEmailRequest{
			Email: hashedUser.Email,
		}
		userUsername := &models.UserByUsernameRequest{
			Username: hashedUser.Username,
		}
		if !service.CheckEmailIsUnique(userEmail) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_EMAIL_ALREADY_USED))
		}
		if !service.CheckUsernameIsUnique(userUsername) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_ALREADY_USED))
		}
		fmt.Println(hashedUser.Username)
		token, err := utils.GenerateJWT(hashedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_TOKEN))
		}

		_, err = service.InsertUser(hashedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.AuthenticationSuccessResponse(token))
	}
}
