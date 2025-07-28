package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/handlers/dto"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/iservice"
	"net/http"
)

type CreditCardHandler struct {
	service iservice.CreditCardManager
}

func NewCreditCardHandler(service iservice.CreditCardManager) *CreditCardHandler {
	return &CreditCardHandler{service: service}
}

// GetAllCreditCards godoc
// @Summary Busca todos os cartões de crédito do usuário
// @Tags CreditCard
// @Produce json
// @Security bearerAuth
// @Success 200 {array} domain.CreditCard
// @Failure 500 {object} map[string]string
// @Router /credit-cards [get]
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

// GetByID godoc
// @Summary Busca um cartão de crédito por ID
// @Tags CreditCard
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID do cartão"
// @Success 200 {object} domain.CreditCard
// @Failure 500 {object} map[string]string
// @Router /credit-cards/{id} [get]
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

// CreateCreditCard godoc
// @Summary Cria um novo cartão de crédito
// @Tags CreditCard
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param credit_card body dto.CreditCardDTO true "Dados do cartão de crédito" example({"card_name":"Nubank Platinum","total_limit":5000,"current_limit":5000,"due_date":10})
// @Success 201 {object} domain.CreditCard
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /credit-cards [post]
func (h *CreditCardHandler) CreateCreditCard(c echo.Context) error {
	var req dto.CreditCardDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	cc := req.ToDomain(userID)
	if err := h.service.Create(cc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, cc)
}

// DeleteCreditCard godoc
// @Summary Remove um cartão de crédito
// @Tags CreditCard
// @Produce json
// @Security bearerAuth
// @Param id path string true "ID do cartão"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /credit-cards/{id} [delete]
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
