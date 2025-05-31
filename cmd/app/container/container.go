package container

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"github.com/misalima/my-budget-planner-backend/internal/core/services"
	postgres "github.com/misalima/my-budget-planner-backend/internal/infra/postgres/repository"
)

type Container struct {
	UserManager       iservice.UserManager
	CategoryManager   iservice.CategoryManager
	AuthManager       iservice.AuthManager
	CreditCardManager iservice.CreditCardManager
}

func NewContainer(pool *pgxpool.Pool) *Container {
	userLoader := postgres.NewUserRepository(pool)
	categoryLoader := postgres.NewCategoryRepository(pool)
	authLoader := postgres.NewAuthRepository(pool)
	creditCardLoader := postgres.NewCreditCardRepository(pool)

	return &Container{
		UserManager:       services.NewUserService(userLoader),
		CategoryManager:   services.NewCategoryService(categoryLoader),
		AuthManager:       services.NewAuthService(authLoader, userLoader),
		CreditCardManager: services.NewCreditCardService(creditCardLoader),
	}
}
