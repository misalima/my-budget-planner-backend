package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type RecurringExpenseUpdateDTO struct {
	ID          string   `json:"id"`
	CategoryID  *int     `json:"category_id,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *string  `json:"date,omitempty"`
	CardID      *string  `json:"card_id,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	Frequency   *string  `json:"frequency,omitempty"`
}

func (dto *RecurringExpenseUpdateDTO) ToDomain(userID uuid.UUID) (domain.RecurringExpense, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return domain.RecurringExpense{}, err
	}
	var date, startDate time.Time
	if dto.Date != nil {
		date, err = time.Parse("2006-01-02", *dto.Date)
		if err != nil {
			return domain.RecurringExpense{}, err
		}
	}
	if dto.StartDate != nil {
		startDate, err = time.Parse("2006-01-02", *dto.StartDate)
		if err != nil {
			return domain.RecurringExpense{}, err
		}
	}
	var endDatePtr *time.Time
	if dto.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *dto.EndDate)
		if err != nil {
			return domain.RecurringExpense{}, err
		}
		endDatePtr = &endDate
	}
	var cardIDPtr *uuid.UUID
	if dto.CardID != nil {
		cardID, err := uuid.Parse(*dto.CardID)
		if err != nil {
			return domain.RecurringExpense{}, err
		}
		cardIDPtr = &cardID
	}
	frequency := ""
	if dto.Frequency != nil {
		frequency = *dto.Frequency
	}
	return domain.RecurringExpense{
		ID:          id,
		UserID:      userID,
		CategoryID:  GetIntValue(dto.CategoryID),
		Amount:      GetFloatValue(dto.Amount),
		Description: dto.Description,
		Date:        date,
		CardID:      cardIDPtr,
		StartDate:   startDate,
		EndDate:     endDatePtr,
		Frequency:   frequency,
	}, nil
}
