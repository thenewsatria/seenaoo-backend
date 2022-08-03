package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
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

func UserLogin(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCredential := &models.Authentication{}
		if err := c.BodyParser(userCredential); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_BODY_PARSER_ERROR_MESSAGE))
		}

		err := validator.IsEmail(userCredential.Credential)
		var loggedUser *models.User

		if err != nil {
			userUsername := &models.UserByUsernameRequest{
				Username: userCredential.Credential,
			}

			loggedUser, err = service.FetchUSerByUsername(userUsername)
			if err != nil {
				c.Status(http.StatusForbidden)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_USERNAME_INCORRECT))
			}
		} else {
			userEmail := &models.UserByEmailRequest{
				Email: userCredential.Credential,
			}
			loggedUser, err = service.FetchUserByEmail(userEmail)
			if err != nil {
				c.Status(http.StatusForbidden)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_EMAIL_INCORRECT))
			}
		}

		checkPass := utils.CheckPasswordHash(userCredential.Password, loggedUser.Password)
		if !checkPass {
			c.Status(http.StatusForbidden)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_PASSWORD_INCORRECT))
		}

		token, err := utils.GenerateJWT(loggedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_TOKEN))
		}
		c.Status(http.StatusOK)
		return c.JSON(presenters.AuthenticationSuccessResponse(token))
	}
}
