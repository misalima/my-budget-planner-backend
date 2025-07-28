package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
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

func (h *CreditCardExpenseHandler) CreateCreditCardExpense(ctx echo.Context) error {
	var dtoReq dto.CreditCardExpenseDTO
	if err := ctx.Bind(&dtoReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
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
			filters.InstallmentsNumber = &n
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

func (h *CreditCardExpenseHandler) GenerateInstallments(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}
	var expense domain.CreditCardExpense
	if err := ctx.Bind(&expense); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	expense.UserID = userID
	if expense.InstallmentsNumber == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "installments number must be greater than zero"})
	}
	installments, err := h.svc.GenerateInstallments(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, installments)
}
