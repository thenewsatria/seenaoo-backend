package presenters

import "github.com/gofiber/fiber/v2"

func ErrorResponse(msg string) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  msg,
	}
}
