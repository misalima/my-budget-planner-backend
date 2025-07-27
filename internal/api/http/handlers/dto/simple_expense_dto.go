package dto

// DTO para SimpleExpense
import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"time"
)

type SimpleExpenseDTO struct {
	UserID      string  `json:"user_id"`
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (dto *SimpleExpenseDTO) ToDomain() (domain.SimpleExpense, error) {
	userID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return domain.SimpleExpense{}, err
	}
	date, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return domain.SimpleExpense{}, err
	}
	return domain.SimpleExpense{
		UserID:      userID,
		CategoryID:  dto.CategoryID,
		Amount:      dto.Amount,
		Description: &dto.Description,
		Date:        date,
	}, nil
}

type SimpleExpenseUpdateDTO struct {
	ID          string   `json:"id"`
	CategoryID  *int     `json:"category_id,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *string  `json:"date,omitempty"`
}

func (dto *SimpleExpenseUpdateDTO) ToDomain(userID uuid.UUID) (domain.SimpleExpense, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return domain.SimpleExpense{}, err
	}
	var date time.Time
	if dto.Date != nil {
		date, err = time.Parse("2006-01-02", *dto.Date)
		if err != nil {
			return domain.SimpleExpense{}, err
		}
	}
	return domain.SimpleExpense{
		ID:          id,
		UserID:      userID,
		CategoryID:  GetIntValue(dto.CategoryID),
		Amount:      GetFloatValue(dto.Amount),
		Description: dto.Description,
		Date:        date,
	}, nil
}

func getIntValue(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func getFloatValue(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
