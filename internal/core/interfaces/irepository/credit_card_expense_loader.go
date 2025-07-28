package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type CreditCardExpenseLoader interface {
	InsertCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error)
	UpdateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error)
	DeleteCreditCardExpense(ctx context.Context, id uuid.UUID) error
	FindCreditCardExpenseByID(ctx context.Context, id uuid.UUID) (domain.CreditCardExpense, error)
	FindCreditCardExpenses(ctx context.Context, userID uuid.UUID, filters CreditCardExpenseFilters) ([]domain.CreditCardExpense, error)
	FindCreditCardExpensesByUser(ctx context.Context, userID uuid.UUID) ([]domain.CreditCardExpense, error)
	FindCreditCardExpensesByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]domain.CreditCardExpense, error)
	InsertInstallments(ctx context.Context, installments []domain.CreditCardExpense) error
}

type CreditCardExpenseFilters struct {
	CategoryID           *int
	CardID               *uuid.UUID
	StartDate            *time.Time
	EndDate              *time.Time
	MinAmount            *float64
	MaxAmount            *float64
	InstallmentsQuantity *int
	ParcelNumber         *int
	Limit                *int
	Offset               *int
}
