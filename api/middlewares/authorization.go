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
	fmt.Println("masuk authorize chech")
	return func(c *fiber.Ctx) error {
		bearerToken, exist := c.GetReqHeaders()["Authorization"]
		if !exist {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_NOT_EXIST))
		}
		token := strings.Split(bearerToken, " ")
		if len(token) < 2 || token[0] != "Bearer" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_ERROR))
		}
		tokenStr := token[1]
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if utils.IsJWTExpired(err) {
				// refresh token logic here
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_EXPIRED))
			}
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_ERROR))
		}

		userUsername := &models.UserByUsernameRequest{
			Username: claims.Username,
		}
		currentUser, err := service.FetchUSerByUsername(userUsername)
		if err != nil {
			c.Status(http.StatusNotFound)
			c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND))
		}
		c.Locals("currentUser", currentUser)
		return c.Next()
	}
}
