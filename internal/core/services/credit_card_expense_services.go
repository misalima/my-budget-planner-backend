package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type CreditCardExpenseService struct {
	repo irepository.CreditCardExpenseLoader
}

func NewCreditCardExpenseService(repo irepository.CreditCardExpenseLoader) *CreditCardExpenseService {
	return &CreditCardExpenseService{repo: repo}
}

func (s *CreditCardExpenseService) CreateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error) {
	return s.repo.InsertCreditCardExpense(ctx, expense)
}

func (s *CreditCardExpenseService) UpdateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error) {
	return s.repo.UpdateCreditCardExpense(ctx, expense)
}

func (s *CreditCardExpenseService) DeleteCreditCardExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	exp, err := s.repo.FindCreditCardExpenseByID(ctx, id)
	if err != nil {
		return err
	}
	if exp.UserID != userID {
		return domain.ErrNotFound
	}
	return s.repo.DeleteCreditCardExpense(ctx, id)
}

func (s *CreditCardExpenseService) GetCreditCardExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.CreditCardExpense, error) {
	exp, err := s.repo.FindCreditCardExpenseByID(ctx, id)
	if err != nil {
		return exp, err
	}
	if exp.UserID != userID {
		return exp, domain.ErrNotFound
	}
	return exp, nil
}

func (s *CreditCardExpenseService) ListCreditCardExpenses(ctx context.Context, userID uuid.UUID, filters irepository.CreditCardExpenseFilters) ([]domain.CreditCardExpense, error) {
	return s.repo.FindCreditCardExpenses(ctx, userID, filters)
}

func (s *CreditCardExpenseService) GenerateInstallments(ctx context.Context, expense domain.CreditCardExpense) ([]domain.CreditCardExpense, error) {
	// Lógica de geração de parcelas pode ser implementada conforme a regra de negócio
	return nil, nil
}

func (s *CreditCardExpenseService) GetCreditCardExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.CreditCardExpenseSummary, error) {
	expenses, err := s.repo.FindCreditCardExpensesByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return domain.CreditCardExpenseSummary{}, err
	}
	var summary domain.CreditCardExpenseSummary
	summary.ByCard = make(map[uuid.UUID]float64)
	summary.ByCategory = make(map[int]float64)
	summary.ByInstallmentsNumber = make(map[int]float64)
	for _, e := range expenses {
		summary.TotalAmount += e.Amount
		summary.TotalCount++
		summary.ByCard[e.CardID] += e.Amount
		summary.ByCategory[e.CategoryID] += e.Amount
		summary.ByInstallmentsNumber[e.InstallmentsNumber] += e.Amount
	}
	if summary.TotalCount > 0 {
		summary.AverageAmount = summary.TotalAmount / float64(summary.TotalCount)
	}
	return summary, nil
}
