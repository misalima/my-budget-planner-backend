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

type SimpleExpenseHandler struct {
	svc iservice.SimpleExpenseManager
}

func NewSimpleExpenseHandler(svc iservice.SimpleExpenseManager) *SimpleExpenseHandler {
	return &SimpleExpenseHandler{svc: svc}
}

func (h *SimpleExpenseHandler) CreateSimpleExpense(ctx echo.Context) error {
	var dtoReq dto.SimpleExpenseDTO
	if err := ctx.Bind(&dtoReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	expense, err := dtoReq.ToDomain()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}
	created, err := h.svc.CreateSimpleExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, created)
}

func (h *SimpleExpenseHandler) GetSimpleExpenseByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}

	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	expense, err := h.svc.GetSimpleExpenseByID(ctx.Request().Context(), id, userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, expense)
}

func (h *SimpleExpenseHandler) ListSimpleExpenses(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id from token"})
	}

	filters := irepository.SimpleExpenseFilters{}

	if v := ctx.QueryParam("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filters.CategoryID = &id
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

	expenses, err := h.svc.ListSimpleExpenses(ctx.Request().Context(), userID, filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, expenses)
}

func (h *SimpleExpenseHandler) UpdateSimpleExpense(ctx echo.Context) error {
	var dtoReq dto.SimpleExpenseUpdateDTO
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
	updated, err := h.svc.UpdateSimpleExpense(ctx.Request().Context(), expense)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, updated)
}

func (h *SimpleExpenseHandler) DeleteSimpleExpense(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
	}
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}
	if err := h.svc.DeleteSimpleExpense(ctx.Request().Context(), id, userID); err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *SimpleExpenseHandler) GetSimpleExpenseSummary(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
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
	summary, err := h.svc.GetSimpleExpenseSummary(ctx.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, summary)
}
