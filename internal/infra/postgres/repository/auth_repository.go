package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"

	"time"
)

type AuthRepository struct {
	Conn *pgxpool.Pool
}

func NewAuthRepository(Conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{Conn: Conn}
}

func (a *AuthRepository) StoreRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string) error {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	createdAt := time.Now()

	sql := `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := a.Conn.Exec(ctx, sql, userId, refreshToken, expirationTime, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepository) GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error) {

	var refreshToken domain.RefreshToken

	sql := `SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`
	err := a.Conn.QueryRow(ctx, sql, token).Scan(&refreshToken)
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (a *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	sql := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := a.Conn.Exec(ctx, sql, token)
	if err != nil {
		return err
	}

	return nil
}
