package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

type User struct {
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
}

type UserDetail struct {
	Username string      `bson:"username" json:"username"`
	Email    string      `bson:"email" json:"email"`
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
		Email:    u.Email,
		Profile:  *userProfile,
	}

	return &fiber.Map{
		"status": "success",
		"data":   userDetail,
	}
}

func UserInsertSuccessResponse(u *models.User) *fiber.Map {
	user := &User{
		Username: u.Username,
	}
	return &fiber.Map{
		"status": "success",
		"data":   user,
	}
}
