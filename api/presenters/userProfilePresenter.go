package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

type UserProfile struct {
	DisplayName     string `bson:"display_name" json:"displayName"`
	AvatarImagePath string `bson:"avatar_image_path" json:"avatarImagePath"`
	BannerImagePath string `bson:"banner_image_path" json:"bannerImagePath"`
	Biography       string `bson:"biography" json:"biography"`
	IsVerified      bool   `bson:"is_verified" json:"isVerified"`
}

type UserProfileDetail struct {
	DisplayName     string `bson:"display_name" json:"displayName" validate:"required,min=1,max=64"`
	AvatarImagePath string `bson:"avatar_image_path" json:"avatarImagePath" validate:"required"`
	BannerImagePath string `bson:"banner_image_path" json:"bannerImagePath" validate:"required"`
	Biography       string `bson:"biography" json:"biography" validate:"omitempty,min=5,max=255"`
	IsVerified      bool   `bson:"is_verified" json:"isVerified" validate:"boolean"`
	Owner           *User  `bson:"owner" json:"owner"`
}

func UserProfileSuccessResponse(up *models.UserProfile) *fiber.Map {
	userProfile := &UserProfile{
		DisplayName:     up.DisplayName,
		AvatarImagePath: up.AvatarImagePath,
		BannerImagePath: up.BannerImagePath,
		Biography:       up.Biography,
		IsVerified:      up.IsVerified,
	}
	return &fiber.Map{
		"status": "success",
		"data":   userProfile,
	}
}

func UserProfileDetailSuccessResponse(up *models.UserProfile, u *models.User) *fiber.Map {
	user := &User{
		Username: u.Username,
		Email:    u.Email,
	}

	userProfile := &UserProfileDetail{
		DisplayName:     up.DisplayName,
		AvatarImagePath: up.AvatarImagePath,
		BannerImagePath: up.BannerImagePath,
		Biography:       up.Biography,
		IsVerified:      up.IsVerified,
		Owner:           user,
	}

	return &fiber.Map{
		"status": "success",
		"data":   userProfile,
	}
}
