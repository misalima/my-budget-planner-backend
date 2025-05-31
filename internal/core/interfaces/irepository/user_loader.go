package irepository

import (
	"context"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type UserLoader interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
