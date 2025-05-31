package router

import (
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/auth"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers"
	"net/http"
)

func LoadRoutes(e *echo.Echo, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, categoryHandler *handlers.CategoryHandler, creditCardHandler *handlers.CreditCardHandler) {
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
	categoryGroup.GET("/:user_id", categoryHandler.GetCategoriesByUserID)
	categoryGroup.POST("", categoryHandler.CreateCategory)
	categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)

	//credit card routes
	creditCardGroup := e.Group("/credit-cards")
	creditCardGroup.Use(auth.JWTMiddleware())
	creditCardGroup.GET("", creditCardHandler.GetAllCreditCards)
	creditCardGroup.POST("", creditCardHandler.CreateCreditCard)
	creditCardGroup.DELETE("/:id", creditCardHandler.DeleteCreditCard)
}
