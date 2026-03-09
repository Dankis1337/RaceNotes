package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levferril/racenotes/backend/internal/models"
	"github.com/levferril/racenotes/backend/internal/services"
)

type SetupHandler struct {
	Service *services.SetupService
}

func NewSetupHandler(service *services.SetupService) *SetupHandler {
	return &SetupHandler{Service: service}
}

type createSetupRequest struct {
	Name                  string  `json:"name" binding:"required"`
	Photo                 *string `json:"photo"`
	BikeName              string  `json:"bike_name" binding:"required"`
	Tires                 string  `json:"tires" binding:"required"`
	ComponentsDescription *string `json:"components_description"`
}

func (h *SetupHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req createSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setup := &models.Setup{
		UserID:                userID,
		Name:                  req.Name,
		Photo:                 req.Photo,
		BikeName:              req.BikeName,
		Tires:                 req.Tires,
		ComponentsDescription: req.ComponentsDescription,
	}

	if err := h.Service.Create(setup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create setup"})
		return
	}

	c.JSON(http.StatusCreated, setup)
}

func (h *SetupHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	setups, err := h.Service.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch setups"})
		return
	}

	c.JSON(http.StatusOK, setups)
}

func (h *SetupHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	setupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid setup id"})
		return
	}

	setup, err := h.Service.GetByID(userID, uint(setupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setup not found"})
		return
	}

	c.JSON(http.StatusOK, setup)
}

func (h *SetupHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	setupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid setup id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(updates, "id")
	delete(updates, "user_id")

	setup, err := h.Service.Update(userID, uint(setupID), updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setup not found"})
		return
	}

	c.JSON(http.StatusOK, setup)
}

func (h *SetupHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	setupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid setup id"})
		return
	}

	if err := h.Service.Delete(userID, uint(setupID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setup not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "setup deleted"})
}
