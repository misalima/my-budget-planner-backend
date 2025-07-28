package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

// CreditCardDTO representa os dados necessários para criar um cartão de crédito.
type CreditCardDTO struct {
	CardName     string  `json:"card_name"`
	TotalLimit   float64 `json:"total_limit"`
	CurrentLimit float64 `json:"current_limit"`
	DueDate      int     `json:"due_date"`
}

// ToDomain converte o DTO para o domínio CreditCard.
func (dto *CreditCardDTO) ToDomain(userID uuid.UUID) *domain.CreditCard {
	return &domain.CreditCard{
		ID:           uuid.New(),
		UserID:       userID,
		CardName:     dto.CardName,
		TotalLimit:   dto.TotalLimit,
		CurrentLimit: dto.CurrentLimit,
		DueDate:      dto.DueDate,
	}
}
