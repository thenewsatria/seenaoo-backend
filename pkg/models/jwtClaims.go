package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	Username string `bson:"username" json:"username"`
	jwt.RegisteredClaims
}
