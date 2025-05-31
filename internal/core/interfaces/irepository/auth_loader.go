package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type AuthLoader interface {
	StoreRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string) error
	GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
