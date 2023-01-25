package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfile struct {
	ID              primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	DisplayName     string             `bson:"display_name" json:"displayName" validate:"required,min=1,max=64"`
	AvatarImagePath string             `bson:"avatar_image_path" json:"avatarImagePath" validate:"required"`
	BannerImagePath string             `bson:"banner_image_path" json:"bannerImagePath" validate:"required"`
	Biography       string             `bson:"biography" json:"biography" validate:"omitempty,min=5,max=255"`
	IsVerified      bool               `bson:"is_verified" json:"isVerified" validate:"boolean"`
	Owner           string             `bson:"owner" json:"owner" validate:"required"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type UserProfileByID struct {
	ID string `bson:"_id" json:"id"`
}

type UserProfileByOwner struct {
	Owner string `bson:"owner" json:"owner"`
}
