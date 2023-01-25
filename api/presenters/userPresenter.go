package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

type User struct {
	Username string `bson:"username" json:"username"`
}

type UserDetail struct {
	Username string      `bson:"username" json:"username"`
	Profile  UserProfile `bson:"profile" json:"profile"`
}

func UserDetailSuccessResponse(u *models.User, p *models.UserProfile) *fiber.Map {
	userProfile := &UserProfile{
		DisplayName:     p.DisplayName,
		AvatarImagePath: p.AvatarImagePath,
		BannerImagePath: p.BannerImagePath,
		Biography:       p.Biography,
		IsVerified:      p.IsVerified,
	}

	userDetail := &UserDetail{
		Username: u.Username,
		Profile:  *userProfile,
	}

	return &fiber.Map{
		"status": true,
		"data":   userDetail,
		"error":  nil,
	}
}

func UserInsertSuccessResponse(u *models.User) *fiber.Map {
	user := &User{
		Username: u.Username,
	}
	return &fiber.Map{
		"status": true,
		"data":   user,
		"error":  nil,
	}
}
