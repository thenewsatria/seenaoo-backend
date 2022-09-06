package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Username        string             `bson:"username" json:"username"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"password"`
	DisplayName     string             `bson:"display_name" json:"displayName"`
	AvatarImagePath string             `bson:"avatar_image_path" json:"avatarImagePath"`
	Biography       string             `bson:"biography" json:"biography"`
	IsVerified      bool               `bson:"is_verified" json:"isVerified"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

type UserByEmailRequest struct {
	Email string `bson:"email" json:"email"`
}

type UserByUsernameRequest struct {
	Username string `bson:"username" json:"username"`
}
