package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Blackmamoth/fileforte/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	Access  TokenType = "ACCESS_TOKEN"
	Refresh TokenType = "REFRESH_TOKEN"
)

func GenerateJWTToken(userId string, r *http.Request, tokenType TokenType) (string, error) {
	expiration := time.Minute * time.Duration(config.JWTConfig.JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userId,
		"remoteAddr": r.RemoteAddr,
		"expiresAt":  expiration,
	})
	switch tokenType {
	case Access:
		tokenStr, err := token.SignedString([]byte(config.JWTConfig.JWT_ACCESS_TOKEN_SECRET))
		if err != nil {
			return "", err
		}
		return tokenStr, nil
	case Refresh:
		tokenStr, err := token.SignedString([]byte(config.JWTConfig.JWT_REFRESH_TOKEN_SECRET))
		if err != nil {
			return "", err
		}
		return tokenStr, nil
	default:
		return "", fmt.Errorf("invalid token type")

	}

}

func ValidateJWTTOken(tokenString string, tokenType TokenType) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		switch tokenType {
		case Access:
			return []byte(config.JWTConfig.JWT_ACCESS_TOKEN_SECRET), nil
		case Refresh:
			return []byte(config.JWTConfig.JWT_REFRESH_TOKEN_SECRET), nil
		default:
			return nil, fmt.Errorf("invalid token type")

		}
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashedPassword(hashed string, plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plainText))
	return err == nil
}
