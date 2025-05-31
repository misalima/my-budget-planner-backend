package domain

import "github.com/google/uuid"

type Category struct {
	ID     int       `json:"id"`
	Name   string    `json:"category_name"`
	UserID uuid.UUID `json:"user_id"`
}
