package iservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type RecurringExpenseManager interface {
	CreateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error)
	UpdateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error)
	DeleteRecurringExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetRecurringExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.RecurringExpense, error)
	ListRecurringExpenses(ctx context.Context, userID uuid.UUID, filters irepository.RecurringExpenseFilters) ([]domain.RecurringExpense, error)
	GenerateRecurringExpenses(ctx context.Context, userID uuid.UUID, targetDate time.Time) error
	GetRecurringExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.RecurringExpenseSummary, error)
}
