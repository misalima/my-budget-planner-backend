package iservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type SimpleExpenseManager interface {
	CreateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error)
	UpdateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error)
	DeleteSimpleExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetSimpleExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.SimpleExpense, error)
	ListSimpleExpenses(ctx context.Context, userID uuid.UUID, filters irepository.SimpleExpenseFilters) ([]domain.SimpleExpense, error)
	GetSimpleExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.SimpleExpenseSummary, error)
}
