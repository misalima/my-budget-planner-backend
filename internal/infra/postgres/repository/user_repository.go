package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"log"
)

// UserRepository is a struct that defines the irepository for the user
type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(Conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{Conn: Conn}
}

func (u UserRepository) CreateUser(ctx context.Context, user *domain.User) error {

	sql := `INSERT INTO users (username, first_name, last_name, email, password_hash) 
			VALUES ($1, $2, $3, $4, $5)`
	log.Print("Executing query")
	tag, err := u.Conn.Exec(ctx, sql, user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows were inserted")
	}

	return nil
}

func (u UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sql := `SELECT "ID", username, first_name, last_name, password_hash, email, profile_picture, income, expenditure_limit, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := u.Conn.QueryRow(ctx, sql, email).Scan(
		&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.Email,
		&user.ProfilePicture, &user.Income, &user.ExpenditureLimit, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
