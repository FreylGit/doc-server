package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"unicode"
)

func HashPasswordT(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func IsValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
		upperCount int
		lowerCount int
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
			upperCount++
		case unicode.IsLower(char):
			hasLower = true
			lowerCount++
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial && upperCount >= 1 && lowerCount >= 1
}

func GenerateAuthToken(userId int64, login string, key []byte) string {
	id := strconv.FormatInt(userId, 10)
	claims := jwt.StandardClaims{
		Id:        id,
		Subject:   login,
		ExpiresAt: time.Now().Add(time.Hour * 50).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(key)

	return tokenStr
}

func ValidateToken(token, expected string) error {
	if token != expected {
		return errors.New("invalid token")
	}
	return nil
}

func ParseToken(tokenStr string, secretKey []byte) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return &jwt.StandardClaims{}, errors.New("invalid token")
	}

	return token.Claims.(*jwt.StandardClaims), nil
}
