package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type RecurringExpenseLoader interface {
	InsertRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error)
	UpdateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error)
	DeleteRecurringExpense(ctx context.Context, id uuid.UUID) error
	FindRecurringExpenseByID(ctx context.Context, id uuid.UUID) (domain.RecurringExpense, error)
	FindRecurringExpenses(ctx context.Context, userID uuid.UUID, filters RecurringExpenseFilters) ([]domain.RecurringExpense, error)
	FindRecurringExpensesByUser(ctx context.Context, userID uuid.UUID) ([]domain.RecurringExpense, error)
	FindRecurringExpensesByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]domain.RecurringExpense, error)
	InsertGeneratedRecurringExpenses(ctx context.Context, expenses []domain.RecurringExpense) error
}

type RecurringExpenseFilters struct {
	CategoryID *int
	CardID     *uuid.UUID
	Frequency  *string
	StartDate  *time.Time
	EndDate    *time.Time
	MinAmount  *float64
	MaxAmount  *float64
	Limit      *int
	Offset     *int
}
