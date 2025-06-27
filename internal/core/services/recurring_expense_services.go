package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type RecurringExpenseService struct {
	repo irepository.RecurringExpenseLoader
}

func NewRecurringExpenseService(repo irepository.RecurringExpenseLoader) *RecurringExpenseService {
	return &RecurringExpenseService{repo: repo}
}

func (s *RecurringExpenseService) CreateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error) {
	return s.repo.InsertRecurringExpense(ctx, expense)
}

func (s *RecurringExpenseService) UpdateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error) {
	return s.repo.UpdateRecurringExpense(ctx, expense)
}

func (s *RecurringExpenseService) DeleteRecurringExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	exp, err := s.repo.FindRecurringExpenseByID(ctx, id)
	if err != nil {
		return err
	}
	if exp.UserID != userID {
		return domain.ErrNotFound
	}
	return s.repo.DeleteRecurringExpense(ctx, id)
}

func (s *RecurringExpenseService) GetRecurringExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.RecurringExpense, error) {
	exp, err := s.repo.FindRecurringExpenseByID(ctx, id)
	if err != nil {
		return exp, err
	}
	if exp.UserID != userID {
		return exp, domain.ErrNotFound
	}
	return exp, nil
}

func (s *RecurringExpenseService) ListRecurringExpenses(ctx context.Context, userID uuid.UUID, filters irepository.RecurringExpenseFilters) ([]domain.RecurringExpense, error) {
	return s.repo.FindRecurringExpenses(ctx, userID, filters)
}

func (s *RecurringExpenseService) GenerateRecurringExpenses(ctx context.Context, userID uuid.UUID, targetDate time.Time) error {
	// Lógica de geração pode ser implementada conforme a regra de negócio
	return nil
}

func (s *RecurringExpenseService) GetRecurringExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.RecurringExpenseSummary, error) {
	expenses, err := s.repo.FindRecurringExpensesByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return domain.RecurringExpenseSummary{}, err
	}
	var summary domain.RecurringExpenseSummary
	summary.ByFrequency = make(map[string]float64)
	summary.ByCategory = make(map[int]float64)
	for _, e := range expenses {
		summary.TotalAmount += e.Amount
		summary.TotalCount++
		summary.ByFrequency[e.Frequency] += e.Amount
		summary.ByCategory[e.CategoryID] += e.Amount
	}
	if summary.TotalCount > 0 {
		summary.AverageAmount = summary.TotalAmount / float64(summary.TotalCount)
	}
	return summary, nil
}
