package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(passw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passw), 14)
	return string(bytes), err
}

func CheckPasswordHash(passw, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passw))
	return err == nil
}
