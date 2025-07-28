package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type SimpleExpenseLoader interface {
	InsertSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error)
	UpdateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error)
	DeleteSimpleExpense(ctx context.Context, expenseId uuid.UUID) error
	FindSimpleExpenseByID(ctx context.Context, expenseId uuid.UUID) (domain.SimpleExpense, error)
	FindSimpleExpenses(ctx context.Context, userId uuid.UUID, filters SimpleExpenseFilters) ([]domain.SimpleExpense, error)
	FindSimpleExpensesByUser(ctx context.Context, userId uuid.UUID) ([]domain.SimpleExpense, error)
	FindSimpleExpensesByDateRange(ctx context.Context, userId uuid.UUID, startDate, endDate time.Time) ([]domain.SimpleExpense, error)
}

type SimpleExpenseFilters struct {
	CategoryID *int
	StartDate  *time.Time
	EndDate    *time.Time
	MinAmount  *float64
	MaxAmount  *float64
	Limit      *int
	Offset     *int
}
