package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

func GenerateJWT(user *models.User) (string, error) {
	jwtKey := os.Getenv("SECRET_KEY")
	jwtExp := os.Getenv("JWT_EXP_TIME")
	expTime, _ := strconv.Atoi(jwtExp)

	claims := &models.JwtClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expTime) * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func ParseJWT(tokenStr string) (*models.JwtClaims, error) {
	jwtKey := os.Getenv("SECRET_KEY")
	claims := &models.JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	} else {
		claims, _ := token.Claims.(*models.JwtClaims)
		return claims, nil
	}
}

func IsJWTExpired(parseError error) bool {
	validationError := parseError.(*jwt.ValidationError)
	return validationError.Errors == jwt.ValidationErrorExpired
}
