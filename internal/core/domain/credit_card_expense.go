package domain

import (
	"github.com/google/uuid"
	"time"
)

type CreditCardExpense struct {
	ID                 uuid.UUID `json:"id" db:"ID"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	CategoryID         int       `json:"category_id" db:"category_id"`
	Amount             float64   `json:"amount" db:"amount"`
	Description        *string   `json:"description" db:"description"`
	Date               time.Time `json:"date" db:"date"`
	CardID             uuid.UUID `json:"card_id" db:"card_id"`
	InstallmentAmount  float64   `json:"installment_amount" db:"installment_amount"`
	InstallmentsNumber int       `json:"installments_number" db:"installments_number"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
