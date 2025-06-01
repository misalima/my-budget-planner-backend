package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type RecurringExpenseRepository struct {
	db *pgxpool.Pool
}

func (r RecurringExpenseRepository) InsertRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error) {
	query := `
		INSERT INTO recurring_expense (user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	now := time.Now()
	expense.CreatedAt = now
	expense.UpdatedAt = now

	err := r.db.QueryRow(ctx, query,
		expense.UserID,
		expense.CategoryID,
		expense.Amount,
		expense.Description,
		expense.Date,
		expense.CardID,
		expense.StartDate,
		expense.EndDate,
		expense.Frequency,
		expense.CreatedAt,
		expense.UpdatedAt,
	).Scan(&expense.ID)

	if err != nil {
		return domain.RecurringExpense{}, fmt.Errorf("failed to insert recurring expense: %w", err)
	}

	return expense, nil
}

func (r RecurringExpenseRepository) UpdateRecurringExpense(ctx context.Context, expense domain.RecurringExpense) (domain.RecurringExpense, error) {
	query := `
		UPDATE recurring_expense 
		SET category_id = $2, amount = $3, description = $4, date = $5, card_id = $6, start_date = $7, end_date = $8, frequency = $9, updated_at = $10
		WHERE id = $1 AND user_id = $11
		RETURNING id, user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at`

	expense.UpdatedAt = time.Now()

	err := r.db.QueryRow(ctx, query,
		expense.ID, expense.CategoryID, expense.Amount, expense.Description, expense.Date,
		expense.CardID, expense.StartDate, expense.EndDate, expense.Frequency, expense.UpdatedAt, expense.UserID,
	).Scan(
		&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
		&expense.Date, &expense.CardID, &expense.StartDate, &expense.EndDate, &expense.Frequency,
		&expense.CreatedAt, &expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.RecurringExpense{}, fmt.Errorf("recurring expense not found or access denied")
		}
		return domain.RecurringExpense{}, fmt.Errorf("failed to update recurring expense: %w", err)
	}

	return expense, nil
}

func (r RecurringExpenseRepository) DeleteRecurringExpense(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recurring_expense WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recurring expense: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("recurring expense not found")
	}

	return nil
}

func (r RecurringExpenseRepository) FindRecurringExpenseByID(ctx context.Context, id uuid.UUID) (domain.RecurringExpense, error) {
	query := `
		SELECT id, user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at
		FROM recurring_expense 
		WHERE id = $1`

	var expense domain.RecurringExpense

	err := r.db.QueryRow(ctx, query, id).Scan(
		&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
		&expense.Date, &expense.CardID, &expense.StartDate, &expense.EndDate, &expense.Frequency,
		&expense.CreatedAt, &expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.RecurringExpense{}, fmt.Errorf("recurring expense not found")
		}
		return domain.RecurringExpense{}, fmt.Errorf("failed to find recurring expense: %w", err)
	}

	return expense, nil
}

func (r RecurringExpenseRepository) FindRecurringExpenses(ctx context.Context, userID uuid.UUID, filters irepository.RecurringExpenseFilters) ([]domain.RecurringExpense, error) {
	query := `
		SELECT id, user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at
		FROM recurring_expense 
		WHERE user_id = $1`

	var args []interface{}
	args = append(args, userID)
	argCount := 2

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
		argCount++
	}

	if filters.CardID != nil {
		query += fmt.Sprintf(" AND card_id = $%d", argCount)
		args = append(args, *filters.CardID)
		argCount++
	}

	if filters.Frequency != nil {
		query += fmt.Sprintf(" AND frequency = $%d", argCount)
		args = append(args, *filters.Frequency)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	if filters.MinAmount != nil {
		query += fmt.Sprintf(" AND amount >= $%d", argCount)
		args = append(args, *filters.MinAmount)
		argCount++
	}

	if filters.MaxAmount != nil {
		query += fmt.Sprintf(" AND amount <= $%d", argCount)
		args = append(args, *filters.MaxAmount)
		argCount++
	}

	query += " ORDER BY start_date DESC, created_at DESC"

	if filters.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *filters.Limit)
		argCount++
	}

	if filters.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, *filters.Offset)
		argCount++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find recurring expenses: %w", err)
	}
	defer rows.Close()

	var expenses []domain.RecurringExpense
	for rows.Next() {
		var expense domain.RecurringExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.StartDate, &expense.EndDate, &expense.Frequency,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recurring expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (r RecurringExpenseRepository) FindRecurringExpensesByUser(ctx context.Context, userID uuid.UUID) ([]domain.RecurringExpense, error) {
	query := `
		SELECT id, user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at
		FROM recurring_expense 
		WHERE user_id = $1
		ORDER BY start_date DESC, created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find recurring expenses by user: %w", err)
	}
	defer rows.Close()

	var expenses []domain.RecurringExpense
	for rows.Next() {
		var expense domain.RecurringExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.StartDate, &expense.EndDate, &expense.Frequency,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recurring expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (r RecurringExpenseRepository) FindRecurringExpensesByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]domain.RecurringExpense, error) {
	query := `
		SELECT id, user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at
		FROM recurring_expense 
		WHERE user_id = $1 AND start_date >= $2 AND (end_date IS NULL OR end_date <= $3)
		ORDER BY start_date DESC, created_at DESC`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to find recurring expenses by date range: %w", err)
	}
	defer rows.Close()

	var expenses []domain.RecurringExpense
	for rows.Next() {
		var expense domain.RecurringExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.StartDate, &expense.EndDate, &expense.Frequency,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recurring expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (r RecurringExpenseRepository) InsertGeneratedRecurringExpenses(ctx context.Context, expenses []domain.RecurringExpense) error {
	if len(expenses) == 0 {
		return nil
	}

	query := `
		INSERT INTO recurring_expense (user_id, category_id, amount, description, date, card_id, start_date, end_date, frequency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	batch := &pgx.Batch{}
	now := time.Now()

	for _, expense := range expenses {
		batch.Queue(query,
			expense.UserID, expense.CategoryID, expense.Amount, expense.Description, expense.Date,
			expense.CardID, expense.StartDate, expense.EndDate, expense.Frequency, now, now,
		)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(expenses); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert generated recurring expense %d: %w", i, err)
		}
	}

	return nil
}

func NewRecurringExpenseRepository(db *pgxpool.Pool) *RecurringExpenseRepository {
	return &RecurringExpenseRepository{
		db: db,
	}
}
