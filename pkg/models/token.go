package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Username     string             `bson:"username" json:"username"`
	RefreshToken string             `bson:"refresh_token" json:"refreshToken"`
	UserAgent    string             `bson:"user_agent" json:"userAgent"`
	ClientIP     string             `bson:"client_ip" json:"clientIP"`
	IsBlocked    bool               `bson:"is_blocked" json:"isBlocked"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

type RefreshTokenByUsersUsernameRequest struct {
	Username string `bson:"username" json:"username"`
}

type RefreshAccessToken struct {
	RefreshToken string `bson:"refresh_token" json:"refreshToken"`
}
