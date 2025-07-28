package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type RecurringExpenseDTO struct {
	UserID      string  `json:"user_id"`
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	CardID      string  `json:"card_id"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Frequency   string  `json:"frequency"`
}

func (dto *RecurringExpenseDTO) ToDomain() (domain.RecurringExpense, error) {
	userID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return domain.RecurringExpense{}, err
	}
	var cardID *uuid.UUID
	if dto.CardID != "" {
		parsed, err := uuid.Parse(dto.CardID)
		if err == nil {
			cardID = &parsed
		}
	}
	date, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return domain.RecurringExpense{}, err
	}
	startDate, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return domain.RecurringExpense{}, err
	}
	var endDate *time.Time
	if dto.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", dto.EndDate)
		if err == nil {
			endDate = &parsed
		}
	}
	return domain.RecurringExpense{
		UserID:      userID,
		CategoryID:  dto.CategoryID,
		Amount:      dto.Amount,
		Description: &dto.Description,
		Date:        date,
		CardID:      cardID,
		StartDate:   startDate,
		EndDate:     endDate,
		Frequency:   dto.Frequency,
	}, nil
}
