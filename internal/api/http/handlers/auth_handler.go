package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
)

type AuthHandler struct {
	AuthService iservice.AuthManager
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authService iservice.AuthManager) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// RefreshTokenHandler refreshes the access token
func (a *AuthHandler) RefreshTokenHandler(ctx echo.Context) error {
	var refreshToken domain.RefreshToken

	//parse the request body, with the token
	if err := ctx.Bind(&refreshToken); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate that the required fields are not empty
	if refreshToken.Token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}

	//extract user id data from the jwt token
	userId, err := uuid.Parse(ctx.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	//call the service
	accessToken, err := a.AuthService.RefreshToken(userId, refreshToken.Token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	//handle the response
	return ctx.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}

func (a *AuthHandler) Login(ctx echo.Context) error {
	var req LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if req.Email == "" || req.Password == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
	}

	accessToken, refreshToken, err := a.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"access_token": accessToken, "refresh_token": refreshToken})
}
