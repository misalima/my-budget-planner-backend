package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":     "my-budget-planner",
		"sub":     userID,
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})
	accessTokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":     "my-budget-planner",
		"sub":     userID,
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days expiration
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Access Denied: authenticate to get access to this functionality"})
		},
	})
}
