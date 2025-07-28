package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"strconv"
	"time"
)

type SimpleExpenseRepository struct {
	db *pgxpool.Pool
}

func (s SimpleExpenseRepository) InsertSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error) {
	query := `
		INSERT INTO simple_expense (user_id, category_id, amount, description, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "ID"`

	now := time.Now()
	expense.CreatedAt = now
	expense.UpdatedAt = now

	err := s.db.QueryRow(ctx, query,
		expense.UserID,
		expense.CategoryID,
		expense.Amount,
		expense.Description,
		expense.Date,
		expense.CreatedAt,
		expense.UpdatedAt,
	).Scan(&expense.ID)

	if err != nil {
		return domain.SimpleExpense{}, fmt.Errorf("failed to insert simple expense: %w", err)
	}

	return expense, nil
}

func (s SimpleExpenseRepository) UpdateSimpleExpense(ctx context.Context, expense domain.SimpleExpense) (domain.SimpleExpense, error) {
	query := "UPDATE simple_expense SET "
	var args []interface{}
	argCount := 1

	if expense.CategoryID != 0 {
		query += "category_id = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.CategoryID)
		argCount++
	}
	if expense.Amount != 0 {
		query += "amount = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.Amount)
		argCount++
	}
	if expense.Description != nil {
		query += "description = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *expense.Description)
		argCount++
	}
	if !expense.Date.IsZero() {
		query += "date = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.Date)
		argCount++
	}
	query += "updated_at = $" + strconv.Itoa(argCount)
	args = append(args, time.Now())
	argCount++

	query += " WHERE \"ID\" = $" + strconv.Itoa(argCount) + " AND user_id = $" + strconv.Itoa(argCount+1)
	args = append(args, expense.ID, expense.UserID)

	query += " RETURNING \"ID\", user_id, category_id, amount, description, date, created_at, updated_at"

	row := s.db.QueryRow(ctx, query, args...)

	err := row.Scan(
		&expense.ID,
		&expense.UserID,
		&expense.CategoryID,
		&expense.Amount,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SimpleExpense{}, fmt.Errorf("simple expense not found or access denied")
		}
		return domain.SimpleExpense{}, fmt.Errorf("failed to update simple expense: %w", err)
	}

	return expense, nil
}

func (s SimpleExpenseRepository) DeleteSimpleExpense(ctx context.Context, expenseId uuid.UUID) error {
	query := `DELETE FROM simple_expense WHERE "ID" = $1`

	result, err := s.db.Exec(ctx, query, expenseId)
	if err != nil {
		return fmt.Errorf("failed to delete simple expense: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("simple expense not found")
	}

	return nil
}

func (s SimpleExpenseRepository) FindSimpleExpenseByID(ctx context.Context, expenseId uuid.UUID) (domain.SimpleExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, created_at, updated_at
		FROM simple_expense 
		WHERE "ID" = $1`

	var expense domain.SimpleExpense

	err := s.db.QueryRow(ctx, query, expenseId).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.CategoryID,
		&expense.Amount,
		&expense.Description,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SimpleExpense{}, fmt.Errorf("simple expense not found")
		}
		return domain.SimpleExpense{}, fmt.Errorf("failed to find simple expense: %w", err)
	}

	return expense, nil
}

func (s SimpleExpenseRepository) FindSimpleExpenses(ctx context.Context, userId uuid.UUID, filters irepository.SimpleExpenseFilters) ([]domain.SimpleExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, created_at, updated_at
		FROM simple_expense 
		WHERE user_id = $1`

	var args []interface{}
	args = append(args, userId)
	argCount := 2

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND date >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND date <= $%d", argCount)
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

	query += " ORDER BY date DESC, created_at DESC"

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

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find simple expenses: %w", err)
	}
	defer rows.Close()

	var expenses []domain.SimpleExpense
	for rows.Next() {
		var expense domain.SimpleExpense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.CategoryID,
			&expense.Amount,
			&expense.Description,
			&expense.Date,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan simple expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (s SimpleExpenseRepository) FindSimpleExpensesByUser(ctx context.Context, userId uuid.UUID) ([]domain.SimpleExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, created_at, updated_at
		FROM simple_expense 
		WHERE user_id = $1
		ORDER BY date DESC, created_at DESC`

	rows, err := s.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find simple expenses by user: %w", err)
	}
	defer rows.Close()

	var expenses []domain.SimpleExpense
	for rows.Next() {
		var expense domain.SimpleExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount,
			&expense.Description, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan simple expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (s SimpleExpenseRepository) FindSimpleExpensesByDateRange(ctx context.Context, userId uuid.UUID, startDate, endDate time.Time) ([]domain.SimpleExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, created_at, updated_at
		FROM simple_expense 
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC, created_at DESC`

	rows, err := s.db.Query(ctx, query, userId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to find simple expenses by date range: %w", err)
	}
	defer rows.Close()

	var expenses []domain.SimpleExpense
	for rows.Next() {
		var expense domain.SimpleExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount,
			&expense.Description, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan simple expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func NewSimpleExpenseRepository(db *pgxpool.Pool) *SimpleExpenseRepository {
	return &SimpleExpenseRepository{
		db: db,
	}
}
