package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/misalima/my-budget-planner-backend/cmd/app/container"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/router"
	"github.com/misalima/my-budget-planner-backend/internal/infra/postgres"
	"os"
)

//	@title			My Budget Planner API
//	@version		1.0
//	@description	This is the backend API for My Budget Planner, an expense tracker and budget manager.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Misael Lima
//	@contact.url	http://www.linkedin.com/misaellima
//	@contact.email	misael.alisson14@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/api

//	@securityDefinitions.apikey	bearerAuth
//	@type 						http
// @scheme 						bearer
//	@in							header
//	@name						Authorization
//	@description				Token JWT para autenticação via header Authorization
//	@bearerFormat				JWT

func main() {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // Allow any localhost port
		AllowCredentials: true,          // Allows cookies and other credentials
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
	}))
	connStr, err := loadEnv()
	if err != nil {
		e.Logger.Fatal(err)
	}

	pool, err := postgres.ConnectDB(connStr)
	if err != nil {
		e.Logger.Fatal(err)
	} else {
		e.Logger.Print("Connected to the database")
	}
	defer pool.Close()

	ctn := container.NewContainer(pool)

	setUpHandlers(e, ctn)

	e.Logger.Fatal(e.Start(":8000"))

}

func loadEnv() (string, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return "", errors.New("Error loading .env file")
	}

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s host=%s sslmode=disable",
		os.Getenv("MBP_PG_USER"),
		os.Getenv("MBP_PG_NAME"),
		os.Getenv("MBP_PG_PASSWORD"),
		os.Getenv("MBP_PG_PORT"),
		os.Getenv("MBP_PG_HOST"),
	)

	return connStr, nil
}

func setUpHandlers(e *echo.Echo, container *container.Container) {

	userHandler := handlers.NewUserHandler(container.UserManager)
	authHandler := handlers.NewAuthHandler(container.AuthManager)
	categoryHandler := handlers.NewCategoryHandler(container.CategoryManager)
	creditCardHandler := handlers.NewCreditCardHandler(container.CreditCardManager)
	simpleExpenseHandler := handlers.NewSimpleExpenseHandler(container.ExpenseManagers.SimpleExpenseManager)
	recurringExpenseHandler := handlers.NewRecurringExpenseHandler(container.ExpenseManagers.RecurringExpenseManager)
	creditCardExpenseHandler := handlers.NewCreditCardExpenseHandler(container.ExpenseManagers.CreditCardExpenseManager)

	router.LoadRoutes(e, userHandler, authHandler, categoryHandler, creditCardHandler, simpleExpenseHandler, recurringExpenseHandler, creditCardExpenseHandler)
}
