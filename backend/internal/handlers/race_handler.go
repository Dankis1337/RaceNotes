package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levferril/racenotes/backend/internal/models"
	"github.com/levferril/racenotes/backend/internal/services"
)

type RaceHandler struct {
	Service *services.RaceService
}

func NewRaceHandler(service *services.RaceService) *RaceHandler {
	return &RaceHandler{Service: service}
}

type createRaceRequest struct {
	Name              string   `json:"name" binding:"required"`
	Date              string   `json:"date" binding:"required"`
	Type              string   `json:"type" binding:"required"`
	Photo             *string  `json:"photo"`
	SetupID           *uint    `json:"setup_id"`
	BikeName          *string  `json:"bike_name"`
	Tires             *string  `json:"tires"`
	TirePressureFront *float64 `json:"tire_pressure_front"`
	TirePressureRear  *float64 `json:"tire_pressure_rear"`
	OtherComponents   *string  `json:"other_components"`
	Temperature       *int     `json:"temperature"`
	Conditions        *string  `json:"conditions"`
	Wind              *string  `json:"wind"`
	RoadConditions    *string  `json:"road_conditions"`
	NutritionPlan     *string  `json:"nutrition_plan"`
	Result            *string  `json:"result"`
	Rating            *int     `json:"rating"`
	Feelings          *string  `json:"feelings"`
	IsCompleted       bool     `json:"is_completed"`
}

func (h *RaceHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req createRaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SetupID == nil && (req.BikeName == nil || req.Tires == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bike_name and tires are required when setup_id is not provided"})
		return
	}

	race := &models.Race{
		UserID:            userID,
		Name:              req.Name,
		Date:              req.Date,
		Type:              req.Type,
		Photo:             req.Photo,
		SetupID:           req.SetupID,
		BikeName:          req.BikeName,
		Tires:             req.Tires,
		TirePressureFront: req.TirePressureFront,
		TirePressureRear:  req.TirePressureRear,
		OtherComponents:   req.OtherComponents,
		Temperature:       req.Temperature,
		Conditions:        req.Conditions,
		Wind:              req.Wind,
		RoadConditions:    req.RoadConditions,
		NutritionPlan:     req.NutritionPlan,
		Result:            req.Result,
		Rating:            req.Rating,
		Feelings:          req.Feelings,
		IsCompleted:       req.IsCompleted,
	}

	if err := h.Service.Create(race); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create race"})
		return
	}

	c.JSON(http.StatusCreated, race)
}

func (h *RaceHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	raceType := c.Query("type")
	isCompleted := c.Query("is_completed")
	setupID := c.Query("setup_id")

	races, err := h.Service.List(userID, raceType, isCompleted, setupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch races"})
		return
	}

	c.JSON(http.StatusOK, races)
}

func (h *RaceHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	raceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
		return
	}

	race, err := h.Service.GetByID(userID, uint(raceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
		return
	}

	c.JSON(http.StatusOK, race)
}

func (h *RaceHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	raceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(updates, "id")
	delete(updates, "user_id")

	race, err := h.Service.Update(userID, uint(raceID), updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
		return
	}

	c.JSON(http.StatusOK, race)
}

func (h *RaceHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	raceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid race id"})
		return
	}

	if err := h.Service.Delete(userID, uint(raceID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "race not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "race deleted"})
}
