package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryService iservice.CategoryManager
}

func NewCategoryHandler(categoryService iservice.CategoryManager) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory godoc
// @Summary Cria uma nova categoria
// @Tags Category
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param category body dto.CreateCategoryDTO true "Dados da categoria" example({"name":"Alimentação"})
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category [post]
func (h *CategoryHandler) CreateCategory(ctx echo.Context) error {
	var req dto.CreateCategoryDTO
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "category name is required"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user id is required"})
	}
	category := req.ToDomain(userID)
	err := h.categoryService.CreateCategory(category)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return ctx.JSON(http.StatusCreated, category)
}

// GetCategoriesByUserID godoc
// @Summary Lista categorias do usuário autenticado
// @Tags Category
// @Produce json
// @Security bearerAuth
// @Success 200 {array} domain.Category
// @Failure 500 {object} map[string]string
// @Router /category [get]
func (h *CategoryHandler) GetCategoriesByUserID(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}
	categories, err := h.categoryService.GetCategoriesByUserID(userID)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return ctx.JSON(http.StatusOK, categories)
}

// DeleteCategory godoc
// @Summary Remove uma categoria
// @Tags Category
// @Produce json
// @Security bearerAuth
// @Param id path int true "ID da categoria"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/{id} [delete]
func (h *CategoryHandler) DeleteCategory(ctx echo.Context) error {
	categoryId := ctx.Param("id")
	intCategoryId, err := strconv.Atoi(categoryId)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "user id is required"})
	}
	err = h.categoryService.DeleteCategory(intCategoryId)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
