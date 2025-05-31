package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CreditCardLoader interface {
	FetchAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CreditCard, error)
	FetchOneByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error)
	Create(ctx context.Context, cc *domain.CreditCard) error
	Delete(ctx context.Context, id uuid.UUID) error
}
