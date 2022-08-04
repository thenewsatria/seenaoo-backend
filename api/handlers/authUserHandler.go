package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/refreshtokens"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(userService users.Service, refreshTokenService refreshtokens.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hashedUser := c.Locals("hashedUser").(*models.User)
		userEmail := &models.UserByEmailRequest{
			Email: hashedUser.Email,
		}
		userUsername := &models.UserByUsernameRequest{
			Username: hashedUser.Username,
		}
		if !userService.CheckEmailIsUnique(userEmail) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_EMAIL_ALREADY_USED_ERROR_MESSAGE))
		}
		if !userService.CheckUsernameIsUnique(userUsername) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_ALREADY_USED_ERROR_MESSAGE))
		}
		accessTokenStr, err := utils.GenerateAccessToken(hashedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE))
		}

		refreshTokenStr, err := utils.GenerateRefreshToken(hashedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE))
		}

		_, err = userService.InsertUser(hashedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		refreshToken := &models.RefreshToken{
			Username:     hashedUser.Username,
			RefreshToken: refreshTokenStr,
			UserAgent:    string(c.Context().UserAgent()),
			ClientIP:     c.IP(),
		}

		_, err = refreshTokenService.InsertRefreshToken(refreshToken)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.AuthenticationSuccessResponse(accessTokenStr, refreshTokenStr))
	}
}

func UserLogin(userService users.Service, refreshTokenService refreshtokens.Service) fiber.Handler {
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

			loggedUser, err = userService.FetchUSerByUsername(userUsername)
			if err != nil {
				c.Status(http.StatusForbidden)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_USERNAME_INCORRECT_ERROR_MESSAGE))
			}
		} else {
			userEmail := &models.UserByEmailRequest{
				Email: userCredential.Credential,
			}
			loggedUser, err = userService.FetchUserByEmail(userEmail)
			if err != nil {
				c.Status(http.StatusForbidden)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_EMAIL_INCORRECT_ERROR_MESSAGE))
			}
		}

		checkPass := utils.CheckPasswordHash(userCredential.Password, loggedUser.Password)
		if !checkPass {
			c.Status(http.StatusForbidden)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_PASSWORD_INCORRECT_ERROR_MESSAGE))
		}

		accessToken, err := utils.GenerateAccessToken(loggedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE))
		}

		refreshToken, err := utils.GenerateRefreshToken(loggedUser)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE))
		}

		//TODO: Add Update refreshToken based on user's username
		loggedUserUsername := &models.RefreshTokenByUsersUsername{Username: loggedUser.Username}
		userRefToken, err := refreshTokenService.FetchRefreshTokenByUsername(loggedUserUsername)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		userRefToken.RefreshToken = refreshToken
		userRefToken.IsBlocked = false
		userRefToken.ClientIP = c.IP()
		userRefToken.UserAgent = string(c.Context().UserAgent())

		updatedToken, err := refreshTokenService.UpdateRefreshToken(userRefToken)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.AuthenticationSuccessResponse(accessToken, updatedToken.RefreshToken))
	}
}
