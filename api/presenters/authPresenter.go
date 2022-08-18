package presenters

import "github.com/gofiber/fiber/v2"

type Authentication struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func AuthenticationSuccessResponse(accessToken string, refreshToken string) *fiber.Map {
	auth := &Authentication{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &fiber.Map{
		"status": true,
		"data":   auth,
		"error":  nil,
	}
}
