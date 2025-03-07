package router

import (
	"github.com/labstack/echo/v4"
	"my-budget-planner/cmd/app/auth"
	"my-budget-planner/cmd/app/handlers"
	"net/http"
)

func LoadRoutes(e *echo.Echo, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, categoryHandler *handlers.CategoryHandler) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/health", handlers.HealthHandler)
	e.POST("/user", userHandler.CreateUserHandler)

	//auth routes
	e.POST("/auth/login", authHandler.Login)
	e.GET("/auth/refresh", authHandler.RefreshTokenHandler)

	//category routes
	categoryGroup := e.Group("/category")
	categoryGroup.Use(auth.JWTMiddleware())
	categoryGroup.GET("/:user_id", categoryHandler.GetCategoriesByUserID)
	categoryGroup.POST("", categoryHandler.CreateCategory)
	categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
}
