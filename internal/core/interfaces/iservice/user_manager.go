package iservice

import "github.com/misalima/my-budget-planner-backend/internal/core/domain"

type UserManager interface {
	RegisterUser(user *domain.User) error
}
