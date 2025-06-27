package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CreditCardRepository struct {
	db *pgxpool.Pool
}

func NewCreditCardRepository(db *pgxpool.Pool) *CreditCardRepository {
	return &CreditCardRepository{db: db}
}

func (r *CreditCardRepository) FetchAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CreditCard, error) {
	rows, err := r.db.Query(ctx, "SELECT 'ID', user_id, card_name, total_limit, current_limit, due_date, created_at, updated_at FROM credit_cards WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creditCards := []domain.CreditCard{}
	for rows.Next() {
		var cc domain.CreditCard
		if err := rows.Scan(&cc.ID, &cc.UserID, &cc.CardName, &cc.TotalLimit, &cc.CurrentLimit, &cc.DueDate, &cc.CreatedAt, &cc.UpdatedAt); err != nil {
			return nil, err
		}
		creditCards = append(creditCards, cc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return creditCards, nil
}

func (r *CreditCardRepository) FetchOneByID(ctx context.Context, id uuid.UUID) (*domain.CreditCard, error) {
	row := r.db.QueryRow(ctx, "SELECT 'ID', user_id, card_name, total_limit, current_limit, due_date, created_at, updated_at FROM credit_cards WHERE 'ID'=$1", id)

	var cc domain.CreditCard
	if err := row.Scan(&cc.ID, &cc.UserID, &cc.CardName, &cc.TotalLimit, &cc.CurrentLimit, &cc.DueDate, &cc.CreatedAt, &cc.UpdatedAt); err != nil {
		return nil, err
	}
	return &cc, nil
}

func (r *CreditCardRepository) Create(ctx context.Context, cc *domain.CreditCard) error {
	_, err := r.db.Exec(ctx, "INSERT INTO credit_cards ('ID', user_id, card_name, total_limit, current_limit, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		cc.ID, cc.UserID, cc.CardName, cc.TotalLimit, cc.CurrentLimit, cc.DueDate, cc.CreatedAt, cc.UpdatedAt)
	return err
}

func (r *CreditCardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM credit_cards WHERE 'ID'=$1", id)
	return err
}
