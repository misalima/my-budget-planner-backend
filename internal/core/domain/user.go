package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	ProfilePicture   string    `json:"profile_picture"`
	Income           float64   `json:"income"`
	ExpenditureLimit float64   `json:"expenditure_limit"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
