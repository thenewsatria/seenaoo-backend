package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func HashUserPassword() fiber.Handler {
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
		c.Locals("hashedUser", user)
		return c.Next()
	}
}
