package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
	"strconv"
	"time"
)

type CreditCardExpenseHandler struct {
	svc iservice.CreditCardExpenseManager
}

func NewCreditCardExpenseHandler(svc iservice.CreditCardExpenseManager) *CreditCardExpenseHandler {
	return &CreditCardExpenseHandler{svc: svc}
}

// CreateCreditCardExpense godoc
// @Summary Cria uma nova despesa de cartão de crédito
// @Tags CreditCardExpense
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param expense body dto.CreditCardExpenseDTO true "Dados da despesa de cartão de crédito"
// @Success 201 {object} domain.CreditCardExpense
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/credit-card [post]
func (h *CreditCardExpenseHandler) CreateCreditCardExpense(ctx echo.Context) error {
	var dtoReq dto.CreditCardExpenseDTO
	if err := ctx.Bind(&dtoReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	dtoReq.UserID = userID.String()

	expense, err := dtoReq.ToDomain()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	created, err := h.svc.CreateCreditCardExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, created)
}

// GetCreditCardExpenseByID godoc
// @Summary Busca uma despesa de cartão de crédito por ID
// @Tags CreditCardExpense
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID da despesa"
// @Success 200 {object} domain.CreditCardExpense
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /expenses/credit-card/{id} [get]
func (h *CreditCardExpenseHandler) GetCreditCardExpenseByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	expense, err := h.svc.GetCreditCardExpenseByID(ctx.Request().Context(), id, userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, expense)
}

// ListCreditCardExpenses godoc
// @Summary Lista despesas de cartão de crédito do usuário
// @Tags CreditCardExpense
// @Produce json
// @Security bearerAuth
// @Param category_id query int false "ID da categoria"
// @Param card_id query string false "ID do cartão"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param min_amount query number false "Valor mínimo"
// @Param max_amount query number false "Valor máximo"
// @Param installments_number query int false "Número de parcelas"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset"
// @Success 200 {array} domain.CreditCardExpense
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/credit-card [get]
func (h *CreditCardExpenseHandler) ListCreditCardExpenses(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}

	filters := irepository.CreditCardExpenseFilters{}

	if v := ctx.QueryParam("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters.CategoryID = &id
		}
	}
	if v := ctx.QueryParam("card_id"); v != "" {
		if cardID, err := uuid.Parse(v); err == nil {
			filters.CardID = &cardID
		}
	}
	if v := ctx.QueryParam("start_date"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filters.StartDate = &t
		}
	}
	if v := ctx.QueryParam("end_date"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filters.EndDate = &t
		}
	}
	if v := ctx.QueryParam("min_amount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MinAmount = &f
		}
	}
	if v := ctx.QueryParam("max_amount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MaxAmount = &f
		}
	}
	if v := ctx.QueryParam("installments_number"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filters.InstallmentsQuantity = &n
		}
	}
	if v := ctx.QueryParam("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			filters.Limit = &l
		}
	}
	if v := ctx.QueryParam("offset"); v != "" {
		if o, err := strconv.Atoi(v); err == nil {
			filters.Offset = &o
		}
	}

	expenses, err := h.svc.ListCreditCardExpenses(ctx.Request().Context(), userID, filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, expenses)
}

// UpdateCreditCardExpense godoc
// @Summary Atualiza uma despesa de cartão de crédito
// @Tags CreditCardExpense
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param expense body dto.CreditCardExpenseUpdateDTO true "Dados da despesa de cartão de crédito para atualização"
// @Success 200 {object} domain.CreditCardExpense
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/credit-card [put]
func (h *CreditCardExpenseHandler) UpdateCreditCardExpense(ctx echo.Context) error {
	var dtoReq dto.CreditCardExpenseUpdateDTO
	if err := ctx.Bind(&dtoReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	expense, err := dtoReq.ToDomain(userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	updated, err := h.svc.UpdateCreditCardExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, updated)
}

// DeleteCreditCardExpense godoc
// @Summary Remove uma despesa de cartão de crédito
// @Tags CreditCardExpense
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID da despesa"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /expenses/credit-card/{id} [delete]
func (h *CreditCardExpenseHandler) DeleteCreditCardExpense(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	if err := h.svc.DeleteCreditCardExpense(ctx.Request().Context(), id, userID); err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetCreditCardExpenseSummary godoc
// @Summary Resumo das despesas de cartão de crédito
// @Tags CreditCardExpense
// @Produce json
// @Security bearerAuth
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /expenses/credit-card/summary [get]
func (h *CreditCardExpenseHandler) GetCreditCardExpenseSummary(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	startDateStr := ctx.QueryParam("start_date")
	endDateStr := ctx.QueryParam("end_date")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid start date"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid end date"})
	}
	summary, err := h.svc.GetCreditCardExpenseSummary(ctx.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, summary)
}
