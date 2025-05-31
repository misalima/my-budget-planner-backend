package iservice

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

type CategoryManager interface {
	CreateCategory(category *domain.Category) error
	GetCategoriesByUserID(userId uuid.UUID) ([]domain.Category, error)
	DeleteCategory(categoryId int) error
}
