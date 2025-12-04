package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UtilisateurID uint   `json:"utilisateurID"`
	Email         string `json:"email"`
	Telephone     string `json:"telephone"`
	RoleID        uint   `json:"roleID"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func InitializeJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateToken(utilisateurID, roleID uint, email, telephone string) (string, error) {
	claims := JWTClaims{
		UtilisateurID: utilisateurID,
		Email:         email,
		Telephone:     telephone,
		RoleID:        roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claim, nil
	}

	return nil, errors.New("token invalide")
}
