package domain

import (
	"github.com/google/uuid"
	"time"
)

type CreditCard struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	CardName     string    `json:"card_name"`
	TotalLimit   float64   `json:"total_limit"`
	CurrentLimit float64   `json:"current_limit"`
	DueDate      int       `json:"due_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
