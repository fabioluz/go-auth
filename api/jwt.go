package api

import (
	"auth/domain/users"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	UserID    string `json:"userId"`
	UserEmail string `json:"userEmail"`
	jwt.RegisteredClaims
}

func generateToken(appCtx *Server, user *users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AppClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
		},
	})

	return token.SignedString(appCtx.JwtSecret)
}

func parseToken(appCtx *Server, strToken string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return appCtx.JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AppClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
