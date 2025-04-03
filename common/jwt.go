package common

import (
	"antara-api/internal/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
	"time"
)

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	rolesString := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		rolesString[i] = role.Name
	}

	// Общие данные для обоих токенов
	commonClaims := CustomJWTClaims{
		ID:    user.ID,
		Roles: rolesString,
	}

	expirationAccessSecondsStr := os.Getenv("ACCESS_TOKEN_EXPIRATION_SECONDS")
	expirationAccessSecondsInt, err := strconv.Atoi(expirationAccessSecondsStr)
	if err != nil {
		return nil, nil, errors.New("invalid expiration seconds")
	}

	// Access token
	userAccessClaims := commonClaims
	userAccessClaims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationAccessSecondsInt) * time.Second)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userAccessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	expirationRefreshSecondsStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_SECONDS")
	expirationRefreshSecondsInt, err := strconv.Atoi(expirationRefreshSecondsStr)
	if err != nil {
		return nil, nil, errors.New("invalid expiration seconds")
	}

	// Refresh token
	userRefreshClaims := commonClaims
	userRefreshClaims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationRefreshSecondsInt) * time.Second)),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userRefreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWTSignedAccessToken(signedAccessToken string) (*CustomJWTClaims, error) {
	parsedJwtAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJWTClaims); ok && parsedJwtAccessToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid JWT access token")
	}
}

func ParseJWTSignedRefreshToken(signedRefreshToken string) (*CustomJWTClaims, error) {
	parsedJwtRefreshToken, err := jwt.ParseWithClaims(signedRefreshToken, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else if claims, ok := parsedJwtRefreshToken.Claims.(*CustomJWTClaims); ok && parsedJwtRefreshToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid JWT refresh token")
	}
}

func IsClaimExpired(claims *CustomJWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}

func GetRefreshTokenExpirationTime() time.Time {
	expirationSecondsStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_SECONDS")
	expirationSecondsInt, err := strconv.Atoi(expirationSecondsStr)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return time.Now().Add(time.Duration(expirationSecondsInt) * time.Second)
}
