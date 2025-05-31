package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"time"
)

type CategoryService struct {
	repo irepository.CategoryLoader
}

func NewCategoryService(repo irepository.CategoryLoader) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (c *CategoryService) CreateCategory(category *domain.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.repo.CheckUserExists(ctx, category.UserID)
	if err != nil {
		return err
	}
	return c.repo.CreateCategory(ctx, category)
}

func (c *CategoryService) GetCategoriesByUserID(userId uuid.UUID) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.repo.GetCategoryByUserID(ctx, userId)
}

func (c *CategoryService) DeleteCategory(categoryId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.repo.DeleteCategory(ctx, categoryId)
}
