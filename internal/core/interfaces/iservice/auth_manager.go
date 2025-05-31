package iservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type AuthManager interface {
	Login(email, password string) (accessToken string, refreshToken string, err error)
	SaveRefreshToken(userId uuid.UUID, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, userId uuid.UUID, token string) (domain.RefreshToken, error)
	DeleteRefreshToken(token string) error
	RefreshToken(userId uuid.UUID, token string) (newAccessToken string, err error)
}
