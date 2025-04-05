package handlers

import (
	calculationservice "kalc/internal/calculationService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type mapStrStr map[string]string

type calculationHandler struct {
	service calculationservice.CalculationService
}

func NewCalculationHandler(s calculationservice.CalculationService) *calculationHandler {
	return &calculationHandler{service: s}
}

// .
// ГЕТ историю
func (h *calculationHandler) GetCalculations(c echo.Context) error {
	calculationsHistory, err := h.service.GetAllCalculations()

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapStrStr{"error": "Couldn't find calculations"})
	}

	return c.JSON(http.StatusOK, calculationsHistory)
}

// .
// ПОСТ результат вычислений
func (h *calculationHandler) PostCalculation(c echo.Context) error {
	// получаем запрос
	var req calculationservice.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, mapStrStr{"error": "Invalud request"})
	}

	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, mapStrStr{"error": "Couldn't create new calculation"})
	}

	// возвращаем ответ
	return c.JSON(http.StatusCreated, calc)
}

// .
// ПАТЧ выражения по айди
func (h *calculationHandler) PatchCalculation(c echo.Context) error {

	// .
	var req calculationservice.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, mapStrStr{"error": "Invalid request"})
	}

	calc, err := h.service.UpdateCalculation(c.Param("id"), req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, mapStrStr{"error": "Couldn't update the record"})
	}

	// .
	return c.JSON(http.StatusOK, calc)
}

// .
// ДЕЛ выражения по айди
func (h *calculationHandler) DelCalculation(c echo.Context) error {
	var req calculationservice.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, mapStrStr{"error": "Invalid request"})
	}

	// .
	if err := h.service.DeleteCalculation(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, mapStrStr{"error": "Couldnt delete from bd or bad id"})
	}

	// .
	return c.NoContent(http.StatusNoContent)
}
