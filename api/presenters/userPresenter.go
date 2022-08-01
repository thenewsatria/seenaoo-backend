package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

type User struct {
	Username        string `bson:"username" json:"username"`
	DisplayName     string `bson:"display_name" json:"displayName"`
	AvatarImagePath string `bson:"avatar_image_path" json:"avatarImagePath"`
	Biography       string `bson:"biography" json:"biography"`
	IsVerified      bool   `bson:"is_verified" json:"isVerified"`
}

func UserInsertSuccessResponse(u *models.User) *fiber.Map {
	user := &User{
		Username:        u.Username,
		DisplayName:     u.DisplayName,
		AvatarImagePath: u.AvatarImagePath,
		Biography:       u.Biography,
		IsVerified:      u.IsVerified,
	}
	return &fiber.Map{
		"status": true,
		"data":   user,
		"error":  nil,
	}
}
