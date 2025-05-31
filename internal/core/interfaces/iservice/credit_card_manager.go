package iservice

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CreditCardManager interface {
	GetAllByUserID(userID uuid.UUID) ([]domain.CreditCard, error)
	GetByID(id uuid.UUID) (*domain.CreditCard, error)
	Create(cc *domain.CreditCard) error
	Delete(id uuid.UUID) error
}
