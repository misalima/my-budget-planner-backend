package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
)

type CreditCardHandler struct {
	service iservice.CreditCardManager
}

func NewCreditCardHandler(service iservice.CreditCardManager) *CreditCardHandler {
	return &CreditCardHandler{service: service}
}

func (h *CreditCardHandler) GetAllCreditCards(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	creditCards, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, creditCards)
}

func (h *CreditCardHandler) GetCreditCardByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credit card ID"})
	}

	creditCard, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, creditCard)
}

func (h *CreditCardHandler) CreateCreditCard(c echo.Context) error {
	var cc domain.CreditCard
	if err := c.Bind(&cc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	cc.ID = uuid.New()
	if err := h.service.Create(&cc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, cc)
}

func (h *CreditCardHandler) DeleteCreditCard(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credit card ID"})
	}

	if err := h.service.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
