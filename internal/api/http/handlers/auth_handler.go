package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
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

// RefreshTokenHandler godoc
// @Summary Atualiza o token de acesso (refresh token)
// @Tags Auth
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param refreshToken body dto.RefreshTokenDTO true "Token de atualização" example({"token":"<refresh_token_aqui>"})
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [get]
func (a *AuthHandler) RefreshTokenHandler(ctx echo.Context) error {
	var req dto.RefreshTokenDTO

	//parse the request body, with the token
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate that the required fields are not empty
	if req.Token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}

	//extract user id data from the jwt token
	userId, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userId == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	//call the service
	accessToken, err := a.AuthService.RefreshToken(userId, req.Token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	//handle the response
	return ctx.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}

// Login godoc
// @Summary Realiza login do usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciais de login" default({"email":"misael@gmail.com","password":"Misael123@"})
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
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
