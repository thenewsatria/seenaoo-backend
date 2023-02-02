package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
)

func LimitReach() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Status(http.StatusTooManyRequests)
		return c.JSON(presenters.ErrorResponse("Limit Exceeded"))
	}
}
