package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
)

type UserHandler struct {
	UserService iservice.UserManager
}

func NewUserHandler(userService iservice.UserManager) *UserHandler {
	return &UserHandler{UserService: userService}
}

// HealthHandler godoc
// @Summary Health check do servidor
// @Tags Health
// @Produce plain
// @Success 200 {string} string "Server is running"
// @Router /health [get]
func HealthHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Server is running")
}

// CreateUserHandler godoc
// @Summary Cria um novo usuário
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.CreateUserDTO true "JSON com as informações de Login de usuário."
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(ctx echo.Context) error {
	var req dto.CreateUserDTO
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "All fields (username, first_name, last_name, email, and password) are required"})
	}
	user := domain.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}
	if err := h.UserService.RegisterUser(&user); err != nil {
		return ctx.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
