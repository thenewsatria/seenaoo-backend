package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func CheckAuthorize(service users.Service) fiber.Handler {
	fmt.Println("masuk authorize check")
	return func(c *fiber.Ctx) error {
		bearerToken, exist := c.GetReqHeaders()["Authorization"]
		if !exist {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_NOT_EXIST_ERROR_MESSAGE))
		}
		token := strings.Split(bearerToken, " ")
		if len(token) < 2 || token[0] != "Bearer" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_INVALID_ERROR_MESSAGE))
		}
		tokenStr := token[1]
		claims, err := utils.ParseAccessToken(tokenStr)
		if err != nil {
			if utils.IsTokenExpired(err) {
				// refresh token logic here
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_EXPIRED_ERROR_MESSAGE))
			}
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_INVALID_ERROR_MESSAGE))
		}

		userUsername := &models.UserByUsernameRequest{
			Username: claims.Username,
		}
		currentUser, err := service.FetchUSerByUsername(userUsername)
		if err != nil {
			c.Status(http.StatusNotFound)
			c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
		}
		c.Locals("currentUser", currentUser)
		return c.Next()
	}
}
