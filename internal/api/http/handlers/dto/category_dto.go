package dto

import (
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
)

// CreateCategoryDTO representa os dados necessários para criar uma categoria.
type CreateCategoryDTO struct {
	Name string `json:"name"`
}

// ToDomain converte o DTO para o domínio Category.
func (dto *CreateCategoryDTO) ToDomain(userID uuid.UUID) *domain.Category {
	return &domain.Category{
		Name:   dto.Name,
		UserID: userID,
	}
}
