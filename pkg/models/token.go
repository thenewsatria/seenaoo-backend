package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID           primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Username     string             `bson:"username" json:"username" validate:"required,lowercase,alphanum,max=25,min=5"`
	RefreshToken string             `bson:"refresh_token" json:"refreshToken" validate:"required"`
	UserAgent    string             `bson:"user_agent" json:"userAgent" validate:"required"`
	ClientIP     string             `bson:"client_ip" json:"clientIP" validate:"required,ip"`
	IsBlocked    bool               `bson:"is_blocked" json:"isBlocked" validate:"required,boolean"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type RefreshTokenByUsersUsernameRequest struct {
	Username string `bson:"username" json:"username" validate:"required,lowercase,alphanum,max=25,min=5"`
}

type RefreshAccessToken struct {
	RefreshToken string `bson:"refresh_token" json:"refreshToken" validate:"required"`
}
