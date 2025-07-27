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

type CreditCardExpenseRepository struct {
	db *pgxpool.Pool
}

func (c CreditCardExpenseRepository) InsertCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error) {
	query := `
		INSERT INTO credit_card_expense (user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING "ID"`

	now := time.Now()
	expense.CreatedAt = now
	expense.UpdatedAt = now

	err := c.db.QueryRow(ctx, query,
		expense.UserID, expense.CategoryID, expense.Amount, expense.Description, expense.Date,
		expense.CardID, expense.InstallmentAmount, expense.InstallmentsNumber, expense.CreatedAt, expense.UpdatedAt,
	).Scan(&expense.ID)

	if err != nil {
		return domain.CreditCardExpense{}, fmt.Errorf("failed to insert credit card expense: %w", err)
	}

	return expense, nil
}

func (c CreditCardExpenseRepository) UpdateCreditCardExpense(ctx context.Context, expense domain.CreditCardExpense) (domain.CreditCardExpense, error) {
	query := `UPDATE credit_card_expense SET `
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
	if expense.CardID != uuid.Nil {
		query += "card_id = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.CardID)
		argCount++
	}
	if expense.InstallmentAmount != 0 {
		query += "installment_amount = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.InstallmentAmount)
		argCount++
	}
	if expense.InstallmentsNumber != 0 {
		query += "installments_number = $" + strconv.Itoa(argCount) + ", "
		args = append(args, expense.InstallmentsNumber)
		argCount++
	}
	query += "updated_at = $" + strconv.Itoa(argCount)
	args = append(args, time.Now())
	argCount++

	query += " WHERE \"ID\" = $" + strconv.Itoa(argCount) + " AND user_id = $" + strconv.Itoa(argCount+1)
	args = append(args, expense.ID, expense.UserID)

	query += " RETURNING \"ID\", user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at"

	row := c.db.QueryRow(ctx, query, args...)

	err := row.Scan(
		&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
		&expense.Date, &expense.CardID, &expense.InstallmentAmount, &expense.InstallmentsNumber,
		&expense.CreatedAt, &expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.CreditCardExpense{}, fmt.Errorf("credit card expense not found or access denied")
		}
		return domain.CreditCardExpense{}, fmt.Errorf("failed to update credit card expense: %w", err)
	}

	return expense, nil
}

func (c CreditCardExpenseRepository) DeleteCreditCardExpense(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM credit_card_expense WHERE "ID" = $1`

	result, err := c.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete credit card expense: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("credit card expense not found")
	}

	return nil
}

func (c CreditCardExpenseRepository) FindCreditCardExpenseByID(ctx context.Context, id uuid.UUID) (domain.CreditCardExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at
		FROM credit_card_expense 
		WHERE "ID" = $1`

	var expense domain.CreditCardExpense

	err := c.db.QueryRow(ctx, query, id).Scan(
		&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
		&expense.Date, &expense.CardID, &expense.InstallmentAmount, &expense.InstallmentsNumber,
		&expense.CreatedAt, &expense.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.CreditCardExpense{}, fmt.Errorf("credit card expense not found")
		}
		return domain.CreditCardExpense{}, fmt.Errorf("failed to find credit card expense: %w", err)
	}

	return expense, nil
}

func (c CreditCardExpenseRepository) FindCreditCardExpenses(ctx context.Context, userID uuid.UUID, filters irepository.CreditCardExpenseFilters) ([]domain.CreditCardExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at
		FROM credit_card_expense 
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

	if filters.InstallmentsNumber != nil {
		query += fmt.Sprintf(" AND installments_number = $%d", argCount)
		args = append(args, *filters.InstallmentsNumber)
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

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find credit card expenses: %w", err)
	}
	defer rows.Close()

	var expenses []domain.CreditCardExpense
	for rows.Next() {
		var expense domain.CreditCardExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.InstallmentAmount, &expense.InstallmentsNumber,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (c CreditCardExpenseRepository) FindCreditCardExpensesByUser(ctx context.Context, userID uuid.UUID) ([]domain.CreditCardExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at
		FROM credit_card_expense 
		WHERE user_id = $1
		ORDER BY date DESC, created_at DESC`

	rows, err := c.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find credit card expenses by user: %w", err)
	}
	defer rows.Close()

	var expenses []domain.CreditCardExpense
	for rows.Next() {
		var expense domain.CreditCardExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.InstallmentAmount, &expense.InstallmentsNumber,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (c CreditCardExpenseRepository) FindCreditCardExpensesByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]domain.CreditCardExpense, error) {
	query := `
		SELECT "ID", user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at
		FROM credit_card_expense 
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC, created_at DESC`

	rows, err := c.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to find credit card expenses by date range: %w", err)
	}
	defer rows.Close()

	var expenses []domain.CreditCardExpense
	for rows.Next() {
		var expense domain.CreditCardExpense
		err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.Date, &expense.CardID, &expense.InstallmentAmount, &expense.InstallmentsNumber,
			&expense.CreatedAt, &expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit card expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return expenses, nil
}

func (c CreditCardExpenseRepository) InsertInstallments(ctx context.Context, installments []domain.CreditCardExpense) error {
	if len(installments) == 0 {
		return nil
	}

	query := `
		INSERT INTO credit_card_expense (user_id, category_id, amount, description, date, card_id, installment_amount, installments_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	batch := &pgx.Batch{}
	now := time.Now()

	for _, installment := range installments {
		batch.Queue(query,
			installment.UserID, installment.CategoryID, installment.Amount, installment.Description, installment.Date,
			installment.CardID, installment.InstallmentAmount, installment.InstallmentsNumber, now, now,
		)
	}

	results := c.db.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(installments); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert installment %d: %w", i, err)
		}
	}

	return nil
}

func NewCreditCardExpenseRepository(db *pgxpool.Pool) *CreditCardExpenseRepository {
	return &CreditCardExpenseRepository{
		db: db,
	}
}
