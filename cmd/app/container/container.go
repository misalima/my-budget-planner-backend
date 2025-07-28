package container

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"github.com/misalima/my-budget-planner-backend/internal/core/services"
	postgres "github.com/misalima/my-budget-planner-backend/internal/infra/postgres"
)

type Container struct {
	UserManager       iservice.UserManager
	CategoryManager   iservice.CategoryManager
	AuthManager       iservice.AuthManager
	CreditCardManager iservice.CreditCardManager
	ExpenseManagers   ExpenseManagers
}

type ExpenseManagers struct {
	CreditCardExpenseManager iservice.CreditCardExpenseManager
	SimpleExpenseManager     iservice.SimpleExpenseManager
	RecurringExpenseManager  iservice.RecurringExpenseManager
}

func NewContainer(pool *pgxpool.Pool) *Container {
	userLoader := postgres.NewUserRepository(pool)
	categoryLoader := postgres.NewCategoryRepository(pool)
	authLoader := postgres.NewAuthRepository(pool)
	creditCardLoader := postgres.NewCreditCardRepository(pool)
	creditCardExpenseLoader := postgres.NewCreditCardExpenseRepository(pool)
	simpleExpenseLoader := postgres.NewSimpleExpenseRepository(pool)
	recurringExpenseLoader := postgres.NewRecurringExpenseRepository(pool)

	return &Container{
		UserManager:       services.NewUserService(userLoader),
		CategoryManager:   services.NewCategoryService(categoryLoader),
		AuthManager:       services.NewAuthService(authLoader, userLoader),
		CreditCardManager: services.NewCreditCardService(creditCardLoader),
		ExpenseManagers: ExpenseManagers{
			CreditCardExpenseManager: services.NewCreditCardExpenseService(creditCardExpenseLoader),
			SimpleExpenseManager:     services.NewSimpleExpenseService(simpleExpenseLoader),
			RecurringExpenseManager:  services.NewRecurringExpenseService(recurringExpenseLoader),
		},
	}
}
