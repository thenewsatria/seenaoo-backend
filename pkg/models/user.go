package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Username        string             `bson:"username" json:"username" validate:"required,lowercase,alphanum,max=25,min=5"`
	Email           string             `bson:"email" json:"email" validate:"required,email"`
	Password        string             `bson:"password" json:"password" validate:"required,min=8"`
	DisplayName     string             `bson:"display_name" json:"displayName" validate:"required,min=1,max=64"`
	AvatarImagePath string             `bson:"avatar_image_path" json:"avatarImagePath" validate:"required,file"`
	Biography       string             `bson:"biography" json:"biography" validate:"omitempty,min=5,max=255"`
	IsVerified      bool               `bson:"is_verified" json:"isVerified"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type UserByEmailRequest struct {
	Email string `bson:"email" json:"email" validate:"required,email"`
}

type UserByUsernameRequest struct {
	Username string `bson:"username" json:"username" validate:"required,lowercase,alphanum,max=25,min=5"`
}
