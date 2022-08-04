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
		user := &models.User{}
		if err := c.BodyParser(user); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_BODY_PARSER_ERROR_MESSAGE))
		}
		hashedPw, err := utils.HashPassword(user.Password)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_HASH_PASSWORD_ERROR_MESSAGE))
		}

		user.Password = hashedPw
		userEmail := &models.UserByEmailRequest{
			Email: user.Email,
		}
		userUsername := &models.UserByUsernameRequest{
			Username: user.Username,
		}
		if !userService.CheckEmailIsUnique(userEmail) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_EMAIL_ALREADY_USED_ERROR_MESSAGE))
		}
		if !userService.CheckUsernameIsUnique(userUsername) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_ALREADY_USED_ERROR_MESSAGE))
		}
		accessTokenStr, err := utils.GenerateAccessToken(user)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE))
		}

		refreshTokenStr, err := utils.GenerateRefreshToken(user)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE))
		}

		_, err = userService.InsertUser(user)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		refreshToken := &models.RefreshToken{
			Username:     user.Username,
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

func UserLogout(refreshTokenService refreshtokens.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		loggedOutUsername := &models.RefreshTokenByUsersUsername{
			Username: currentUser.Username,
		}
		userRefToken, err := refreshTokenService.FetchRefreshTokenByUsername(loggedOutUsername)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		userRefToken.RefreshToken = "-"
		userRefToken.IsBlocked = true
		userRefToken.ClientIP = c.IP()
		userRefToken.UserAgent = string(c.Context().UserAgent())

		_, err = refreshTokenService.UpdateRefreshToken(userRefToken)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.AuthenticationSuccessResponse("-", "-"))
	}
}
