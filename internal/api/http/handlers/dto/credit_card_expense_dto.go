package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type CreditCardExpenseDTO struct {
	UserID               string  `json:"user_id,omitempty"`
	CategoryID           int     `json:"category_id"`
	Amount               float64 `json:"amount"`
	Description          string  `json:"description"`
	Date                 string  `json:"date"`
	CardID               string  `json:"card_id"`
	InstallmentAmount    float64 `json:"installment_amount"`
	InstallmentsQuantity int     `json:"installments_quantity"`
}

func (dto *CreditCardExpenseDTO) ToDomain() (domain.CreditCardExpense, error) {
	userID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return domain.CreditCardExpense{}, err
	}
	cardID, err := uuid.Parse(dto.CardID)
	if err != nil {
		return domain.CreditCardExpense{}, err
	}
	date, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return domain.CreditCardExpense{}, err
	}
	return domain.CreditCardExpense{
		UserID:               userID,
		CategoryID:           dto.CategoryID,
		Amount:               dto.Amount,
		Description:          &dto.Description,
		Date:                 date,
		CardID:               cardID,
		InstallmentAmount:    dto.InstallmentAmount,
		InstallmentsQuantity: dto.InstallmentsQuantity,
	}, nil
}
