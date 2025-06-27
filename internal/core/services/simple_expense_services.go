package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type SimpleExpenseService struct {
	repo irepository.SimpleExpenseLoader
}

func NewSimpleExpenseService(repo irepository.SimpleExpenseLoader) *SimpleExpenseService {
	return &SimpleExpenseService{repo: repo}
}

func (s *SimpleExpenseService) CreateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error) {
	return s.repo.InsertSimpleExpense(ctx, expense)
}

func (s *SimpleExpenseService) UpdateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error) {
	return s.repo.UpdateSimpleExpense(ctx, expense)
}

func (s *SimpleExpenseService) DeleteSimpleExpense(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	exp, err := s.repo.FindSimpleExpenseByID(ctx, id)
	if err != nil {
		return err
	}
	if exp.UserID != userID {
		return domain.ErrNotFound
	}
	return s.repo.DeleteSimpleExpense(ctx, id)
}

func (s *SimpleExpenseService) GetSimpleExpenseByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.SimpleExpense, error) {
	exp, err := s.repo.FindSimpleExpenseByID(ctx, id)
	if err != nil {
		return exp, err
	}
	if exp.UserID != userID {
		return exp, domain.ErrNotFound
	}
	return exp, nil
}

func (s *SimpleExpenseService) ListSimpleExpenses(ctx context.Context, userID uuid.UUID, filters irepository.SimpleExpenseFilters) ([]domain.SimpleExpense, error) {
	return s.repo.FindSimpleExpenses(ctx, userID, filters)
}

func (s *SimpleExpenseService) GetSimpleExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (domain.SimpleExpenseSummary, error) {
	expenses, err := s.repo.FindSimpleExpensesByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return domain.SimpleExpenseSummary{}, err
	}
	var summary domain.SimpleExpenseSummary
	summary.ByCategory = make(map[int]float64)
	for _, e := range expenses {
		summary.TotalAmount += e.Amount
		summary.TotalCount++
		summary.ByCategory[e.CategoryID] += e.Amount
	}
	if summary.TotalCount > 0 {
		summary.AverageAmount = summary.TotalAmount / float64(summary.TotalCount)
	}
	return summary, nil
}
