package iservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type CreditCardExpenseManager interface {
	CreateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error)
	UpdateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error)
	DeleteCreditCardExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetCreditCardExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.CreditCardExpense, error)
	ListCreditCardExpenses(ctx context.Context, userID uuid.UUID, filters irepository.CreditCardExpenseFilters) ([]domain.CreditCardExpense, error)
	GenerateInstallments(ctx context.Context, expense domain.CreditCardExpense) ([]domain.CreditCardExpense, error)
	GetCreditCardExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.CreditCardExpenseSummary, error)
}
