package router

import (
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/auth"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers"
	"net/http"
)

func LoadRoutes(
	e *echo.Echo,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	creditCardHandler *handlers.CreditCardHandler,
	simpleExpenseHandler *handlers.SimpleExpenseHandler,
	recurringExpenseHandler *handlers.RecurringExpenseHandler,
	creditCardExpenseHandler *handlers.CreditCardExpenseHandler,
) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/health", handlers.HealthHandler)
	e.POST("/signup", userHandler.CreateUserHandler)

	//auth routes
	e.POST("/auth/login", authHandler.Login)
	e.GET("/auth/refresh", authHandler.RefreshTokenHandler)

	//category routes
	categoryGroup := e.Group("/category")
	categoryGroup.Use(auth.JWTMiddleware())
	categoryGroup.Use(auth.ExtractUserIDMiddleware)
	categoryGroup.GET("/:user_id", categoryHandler.GetCategoriesByUserID)
	categoryGroup.POST("", categoryHandler.CreateCategory)
	categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)

	//credit card routes
	creditCardGroup := e.Group("/credit-cards")
	creditCardGroup.Use(auth.JWTMiddleware())
	creditCardGroup.Use(auth.ExtractUserIDMiddleware)
	creditCardGroup.GET("", creditCardHandler.GetAllCreditCards)
	creditCardGroup.POST("", creditCardHandler.CreateCreditCard)
	creditCardGroup.DELETE("/:id", creditCardHandler.DeleteCreditCard)

	// expenses routes
	expenseGroup := e.Group("/expenses")
	expenseGroup.Use(auth.JWTMiddleware())
	expenseGroup.Use(auth.ExtractUserIDMiddleware)

	// Simple expenses
	simpleGroup := expenseGroup.Group("/simple")
	simpleGroup.GET("", simpleExpenseHandler.ListSimpleExpenses)
	simpleGroup.GET("/:id", simpleExpenseHandler.GetSimpleExpenseByID)
	simpleGroup.POST("", simpleExpenseHandler.CreateSimpleExpense)
	simpleGroup.PUT("/", simpleExpenseHandler.UpdateSimpleExpense)
	simpleGroup.DELETE("/:id", simpleExpenseHandler.DeleteSimpleExpense)
	simpleGroup.GET("/summary", simpleExpenseHandler.GetSimpleExpenseSummary)

	// Recurring expenses
	recurringGroup := expenseGroup.Group("/recurring")
	recurringGroup.GET("", recurringExpenseHandler.ListRecurringExpenses)
	recurringGroup.GET("/:id", recurringExpenseHandler.GetRecurringExpenseByID)
	recurringGroup.POST("", recurringExpenseHandler.CreateRecurringExpense)
	recurringGroup.PUT("/", recurringExpenseHandler.UpdateRecurringExpense)
	recurringGroup.DELETE("/:id", recurringExpenseHandler.DeleteRecurringExpense)
	recurringGroup.GET("/summary", recurringExpenseHandler.GetRecurringExpenseSummary)
	recurringGroup.POST("/generate", recurringExpenseHandler.GenerateRecurringExpenses)

	// Credit card expenses
	creditCardExpenseGroup := expenseGroup.Group("/credit-card")
	creditCardExpenseGroup.GET("", creditCardExpenseHandler.ListCreditCardExpenses)
	creditCardExpenseGroup.GET("/:id", creditCardExpenseHandler.GetCreditCardExpenseByID)
	creditCardExpenseGroup.POST("", creditCardExpenseHandler.CreateCreditCardExpense)
	creditCardExpenseGroup.PUT("/", creditCardExpenseHandler.UpdateCreditCardExpense)
	creditCardExpenseGroup.DELETE("/:id", creditCardExpenseHandler.DeleteCreditCardExpense)
	creditCardExpenseGroup.GET("/summary", creditCardExpenseHandler.GetCreditCardExpenseSummary)
	creditCardExpenseGroup.POST("/installments/generate", creditCardExpenseHandler.GenerateInstallments)
}
