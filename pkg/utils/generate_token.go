package utils

import (
	"errors"
	"time"

	"github.com/cabitibaly/configs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UtilisateurID uint   `json:"utilisateurID"`
	Email         string `json:"email"`
	Telephone     string `json:"telephone"`
	RoleID        uint   `json:"roleID"`
	Jti           string `json:"jti"`
	jwt.RegisteredClaims
}

var (
	JwtSecret          []byte
	RefreshTokenSecret []byte
)

func InitializeJWTSecret(secret string, refreshSecret string) {
	JwtSecret = []byte(secret)
	RefreshTokenSecret = []byte(refreshSecret)
}

func GenerateAccessToken(utilisateurID, roleID uint, email, telephone string) (string, error) {
	location, _ := time.LoadLocation("Africa/Ouagadougou")

	expireAt := time.Now().Add(24 * time.Hour).In(location)

	jti := uuid.New().String()

	claims := JWTClaims{
		UtilisateurID: utilisateurID,
		Email:         email,
		Telephone:     telephone,
		RoleID:        roleID,
		Jti:           jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	cacheKey := "access_token:" + jti
	err := configs.SetCache(cacheKey, jti, 24*time.Hour)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claim, nil
	}

	return nil, errors.New("token invalide")
}

func GenerateRefreshToken(utilisateurID, roleID uint, email, telephone string) (string, *time.Time, error) {
	location, _ := time.LoadLocation("Africa/Ouagadougou")

	expireAt := time.Now().Add(30 * 24 * time.Hour).In(location)

	claims := JWTClaims{
		UtilisateurID: utilisateurID,
		Email:         email,
		Telephone:     telephone,
		RoleID:        roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(location)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(RefreshTokenSecret)

	return tokenStr, &expireAt, err
}
