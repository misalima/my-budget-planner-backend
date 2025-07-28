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

type RecurringExpenseHandler struct {
	svc iservice.RecurringExpenseManager
}

func NewRecurringExpenseHandler(svc iservice.RecurringExpenseManager) *RecurringExpenseHandler {
	return &RecurringExpenseHandler{svc: svc}
}

// CreateRecurringExpense godoc
// @Summary Cria uma nova despesa recorrente
// @Tags RecurringExpense
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param expense body dto.RecurringExpenseDTO true "Dados da despesa recorrente"
// @Success 201 {object} domain.RecurringExpense
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/recurring [post]
func (h *RecurringExpenseHandler) CreateRecurringExpense(ctx echo.Context) error {
	var dtoReq dto.RecurringExpenseDTO
	if err := ctx.Bind(&dtoReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	expense, err := dtoReq.ToDomain()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	created, err := h.svc.CreateRecurringExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, created)
}

// GetRecurringExpenseByID godoc
// @Summary Busca uma despesa recorrente por ID
// @Tags RecurringExpense
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID da despesa"
// @Success 200 {object} domain.RecurringExpense
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /expenses/recurring/{id} [get]
func (h *RecurringExpenseHandler) GetRecurringExpenseByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	expense, err := h.svc.GetRecurringExpenseByID(ctx.Request().Context(), id, userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, expense)
}

// ListRecurringExpenses godoc
// @Summary Lista despesas recorrentes do usuário
// @Tags RecurringExpense
// @Produce json
// @Security bearerAuth
// @Param category_id query int false "ID da categoria"
// @Param card_id query string false "ID do cartão"
// @Param frequency query string false "Frequência"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param min_amount query number false "Valor mínimo"
// @Param max_amount query number false "Valor máximo"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset"
// @Success 200 {array} domain.RecurringExpense
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/recurring [get]
func (h *RecurringExpenseHandler) ListRecurringExpenses(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}

	filters := irepository.RecurringExpenseFilters{}

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
	if v := ctx.QueryParam("frequency"); v != "" {
		filters.Frequency = &v
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

	expenses, err := h.svc.ListRecurringExpenses(ctx.Request().Context(), userID, filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, expenses)
}

// UpdateRecurringExpense godoc
// @Summary Atualiza uma despesa recorrente
// @Tags RecurringExpense
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param expense body dto.RecurringExpenseUpdateDTO true "Dados da despesa recorrente para atualização"
// @Success 200 {object} domain.RecurringExpense
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/recurring [put]
func (h *RecurringExpenseHandler) UpdateRecurringExpense(ctx echo.Context) error {
	var dtoReq dto.RecurringExpenseUpdateDTO
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
	updated, err := h.svc.UpdateRecurringExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, updated)
}

// DeleteRecurringExpense godoc
// @Summary Remove uma despesa recorrente
// @Tags RecurringExpense
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID da despesa"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /expenses/recurring/{id} [delete]
func (h *RecurringExpenseHandler) DeleteRecurringExpense(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	if err := h.svc.DeleteRecurringExpense(ctx.Request().Context(), id, userID); err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetRecurringExpenseSummary godoc
// @Summary Resumo das despesas recorrentes
// @Tags RecurringExpense
// @Produce json
// @Security bearerAuth
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/recurring/summary [get]
func (h *RecurringExpenseHandler) GetRecurringExpenseSummary(ctx echo.Context) error {
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
	summary, err := h.svc.GetRecurringExpenseSummary(ctx.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, summary)
}

// GenerateRecurringExpenses godoc
// @Summary Gera despesas recorrentes para uma data alvo
// @Tags RecurringExpense
// @Produce json
// @Security bearerAuth
// @Param target_date query string true "Data alvo (YYYY-MM-DD)"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/recurring/generate [post]
func (h *RecurringExpenseHandler) GenerateRecurringExpenses(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	targetDateStr := ctx.QueryParam("target_date")
	targetDate, err := time.Parse("2006-01-02", targetDateStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid target date"})
	}
	if err := h.svc.GenerateRecurringExpenses(ctx.Request().Context(), userID, targetDate); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"status": "generated"})
}
