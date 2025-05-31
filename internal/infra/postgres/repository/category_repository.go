package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CategoryRepository struct {
	Conn *pgxpool.Pool
}

func NewCategoryRepository(Conn *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{Conn: Conn}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	sql := `INSERT INTO categories (category_name, user_id) VALUES ($1, $2) RETURNING "ID"`

	err := c.Conn.QueryRow(ctx, sql, category.Name, category.UserID).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) GetCategoryByUserID(ctx context.Context, userId uuid.UUID) ([]domain.Category, error) {
	sql := `SELECT * FROM categories WHERE user_id = $1 OR user_id IS NULL`

	rows, err := c.Conn.Query(ctx, sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name, &category.UserID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return categories, nil
}

func (c *CategoryRepository) DeleteCategory(ctx context.Context, categoryId int) error {
	sql := `DELETE FROM categories WHERE "ID" = $1`
	_, err := c.Conn.Exec(ctx, sql, categoryId)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) CheckUserExists(ctx context.Context, id uuid.UUID) error {
	sql := `SELECT COUNT(1) FROM users WHERE "ID" = $1`

	var count int
	err := c.Conn.QueryRow(ctx, sql, id).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}
