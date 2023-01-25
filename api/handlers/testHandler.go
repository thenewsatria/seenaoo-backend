package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/tests"
)

func TestHandler(service tests.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(&fiber.Map{
				"status": false,
				"error":  err.Error(),
			})
		}

		files := form.File["documents"]
		if files == nil {
			c.Status(http.StatusNotFound)
			return c.JSON(&fiber.Map{
				"status": "file Kosong",
			})
		}
		for _, file := range files {
			fmt.Println(file.Filename)
		}

		c.Status(http.StatusInternalServerError)
		return c.JSON(&fiber.Map{
			"status": true,
			"error":  "sabi",
		})
	}
}
