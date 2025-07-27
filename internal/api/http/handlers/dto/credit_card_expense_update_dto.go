package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type CreditCardExpenseUpdateDTO struct {
	ID                 string   `json:"id"`
	CategoryID         *int     `json:"category_id,omitempty"`
	Amount             *float64 `json:"amount,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Date               *string  `json:"date,omitempty"`
	CardID             *string  `json:"card_id,omitempty"`
	InstallmentAmount  *float64 `json:"installment_amount,omitempty"`
	InstallmentsNumber *int     `json:"installments_number,omitempty"`
}

func (dto *CreditCardExpenseUpdateDTO) ToDomain(userID uuid.UUID) (domain.CreditCardExpense, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return domain.CreditCardExpense{}, err
	}
	var date time.Time
	if dto.Date != nil {
		date, err = time.Parse("2006-01-02", *dto.Date)
		if err != nil {
			return domain.CreditCardExpense{}, err
		}
	}
	var cardID uuid.UUID
	if dto.CardID != nil {
		cardID, err = uuid.Parse(*dto.CardID)
		if err != nil {
			return domain.CreditCardExpense{}, err
		}
	}
	return domain.CreditCardExpense{
		ID:                 id,
		UserID:             userID,
		CategoryID:         GetIntValue(dto.CategoryID),
		Amount:             GetFloatValue(dto.Amount),
		Description:        dto.Description,
		Date:               date,
		CardID:             cardID,
		InstallmentAmount:  GetFloatValue(dto.InstallmentAmount),
		InstallmentsNumber: GetIntValue(dto.InstallmentsNumber),
	}, nil
}
