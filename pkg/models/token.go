package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	RefreshToken string             `bson:"refresh_token"`
	UserAgent    string             `bson:"user_agent"`
	ClientIP     string             `bson:"client_ip"`
	IsBlocked    bool               `bson:"is_blocked"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type RefreshTokenByUsersUsername struct {
	Username string `bson:"username"`
}

type RefreshAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}
