package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type CreditCardService struct {
	repo irepository.CreditCardLoader
}

func NewCreditCardService(repo irepository.CreditCardLoader) *CreditCardService {
	return &CreditCardService{repo: repo}
}

func (s *CreditCardService) GetAllByUserID(userID uuid.UUID) ([]domain.CreditCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.FetchAllByUserID(ctx, userID)
}

func (s *CreditCardService) GetByID(id uuid.UUID) (*domain.CreditCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.FetchOneByID(ctx, id)
}

func (s *CreditCardService) Create(cc *domain.CreditCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, cc)
}

func (s *CreditCardService) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.Delete(ctx, id)
}
