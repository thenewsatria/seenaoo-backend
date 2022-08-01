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
