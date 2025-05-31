package irepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CategoryLoader interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByUserID(ctx context.Context, userId uuid.UUID) ([]domain.Category, error)
	DeleteCategory(ctx context.Context, categoryId int) error
	CheckUserExists(ctx context.Context, id uuid.UUID) error
}
