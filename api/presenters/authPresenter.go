package presenters

import "github.com/gofiber/fiber/v2"

type Authentication struct {
	Token string `json:"token"`
}

func AuthenticationSuccessResponse(jwtToken string) *fiber.Map {
	auth := &Authentication{
		Token: jwtToken,
	}
	return &fiber.Map{
		"status": true,
		"data":   auth,
		"error":  nil,
	}
}
