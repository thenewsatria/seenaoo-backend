package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/refreshtokens"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(userService users.Service, refreshTokenService refreshtokens.Service, userProfileService userprofiles.Service) fiber.Handler {
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
		if userService.CheckEmailIsExist(userEmail) {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_EMAIL_ALREADY_USED_ERROR_MESSAGE))
		}
		if userService.CheckUsernameIsExist(userUsername) {
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

		_, err, isValidationError := userService.InsertUser(user)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		userProfile := &models.UserProfile{
			DisplayName:     user.Username,
			AvatarImagePath: "http://localhost:3000/api/v1/static/defaults/default-avatar.png", //default image
			BannerImagePath: "http://localhost:3000/api/v1/static/defaults/default-banner.jpg",
			Owner:           user.Username,
		}

		_, err, isValidationError = userProfileService.InsertProfile(userProfile)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		refreshToken := &models.RefreshToken{
			Username:     user.Username,
			RefreshToken: refreshTokenStr,
			UserAgent:    string(c.Context().UserAgent()),
			ClientIP:     c.IP(),
		}

		_, err, isValidationError = refreshTokenService.InsertRefreshToken(refreshToken)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		c.Status(http.StatusCreated)
		return c.JSON(presenters.AuthenticationSuccessResponse(accessTokenStr, refreshTokenStr))
	}
}

func UserLogin(userService users.Service, refreshTokenService refreshtokens.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCredential := &models.LoginRequest{}
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

			loggedUser, err = userService.FetchUserByUsername(userUsername)
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

		loggedUserUsername := &models.RefreshTokenByUsersUsernameRequest{Username: loggedUser.Username}
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
		userRefToken.ClientIP = c.IP()
		userRefToken.UserAgent = string(c.Context().UserAgent())

		updatedToken, err, isValidationError := refreshTokenService.UpdateRefreshToken(userRefToken)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
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
		loggedOutUsername := &models.RefreshTokenByUsersUsernameRequest{
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
		userRefToken.ClientIP = c.IP()
		userRefToken.UserAgent = string(c.Context().UserAgent())

		_, err, isValidationError := refreshTokenService.UpdateRefreshToken(userRefToken)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.AuthenticationSuccessResponse("-", "-"))
	}
}

func RefreshToken(refreshTokenService refreshtokens.Service, userService users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshRequest := &models.RefreshAccessToken{}
		if err := c.BodyParser(refreshRequest); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_ACCESS_TOKEN_BODY_PARSER_ERROR_MESSAGE))
		}

		claims, err := utils.ParseRefreshoken(refreshRequest.RefreshToken)
		if err != nil {
			if utils.IsTokenExpired(err) {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_REFRESH_TOKEN_EXPIRED_ERROR_MESSAGE))
			}
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_INVALID_ERROR_MESSAGE))
		}
		rtByUname := &models.RefreshTokenByUsersUsernameRequest{
			Username: claims.Username,
		}
		refTok, err := refreshTokenService.FetchRefreshTokenByUsername(rtByUname)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		if refTok.IsBlocked {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_REFRESH_TOKEN_BLOCKED_ERROR_MESSAGE))
		}

		if refTok.RefreshToken == "-" {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_STORED_REFRESH_TOKEN_IS_EMPTY_ERROR_MESSAGE))
		}

		if refTok.RefreshToken != refreshRequest.RefreshToken {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_REFRESH_TOKEN_DIFFERENT_FROM_STORED_ERROR_MESSAGE))
		}

		userByUname := &models.UserByUsernameRequest{
			Username: claims.Username,
		}

		userIssued, err := userService.FetchUserByUsername(userByUname)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
		}

		newAccessToken, err := utils.GenerateAccessToken(userIssued)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE))
		}

		newRefreshToken, err := utils.GenerateRefreshToken(userIssued)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE))
		}

		refTok.RefreshToken = newRefreshToken
		refTok.ClientIP = c.IP()
		refTok.UserAgent = string(c.Context().UserAgent())

		updatedToken, err, isValidationError := refreshTokenService.UpdateRefreshToken(refTok)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.REFRESH_TOKEN_FAIL_TO_UPDATE_STORED_TOKEN_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.AuthenticationSuccessResponse(newAccessToken, updatedToken.RefreshToken))
	}
}
