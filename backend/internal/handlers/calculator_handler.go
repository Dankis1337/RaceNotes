package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levferril/racenotes/backend/internal/services"
)

type CalculatorHandler struct{}

func NewCalculatorHandler() *CalculatorHandler {
	return &CalculatorHandler{}
}

func (h *CalculatorHandler) Calculate(c *gin.Context) {
	var req services.CalculatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := services.CalculateTirePressure(req)
	c.JSON(http.StatusOK, result)
}
