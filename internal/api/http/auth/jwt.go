package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
		"sub":     userID.String(),
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Minute * 60).Unix(),
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
		"sub":     userID.String(),
		"user_id": userID.String(),
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

func ExtractUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token not found"})
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("invalid token type: %T", user)})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims type in *jwt.Token"})
		}

		userIDRaw, ok := claims["user_id"]
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user_id not found in token"})
		}

		userIDStr, ok := userIDRaw.(string)
		if !ok {
			userIDStr = fmt.Sprintf("%v", userIDRaw)
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id format in token"})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
