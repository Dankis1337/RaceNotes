package services

import (
	"fmt"
	"math"
	"strings"
)

type CalculatorRequest struct {
	RiderWeight float64 `json:"rider_weight" binding:"required"`
	BikeWeight  float64 `json:"bike_weight" binding:"required"`
	TireWidth   int     `json:"tire_width" binding:"required"`
	TireType    string  `json:"tire_type" binding:"required"`
	Surface     string  `json:"surface" binding:"required"`
	Conditions  string  `json:"conditions"`
}

type CalculatorResponse struct {
	FrontPressure   float64  `json:"front_pressure"`
	RearPressure    float64  `json:"rear_pressure"`
	Unit            string   `json:"unit"`
	Recommendations []string `json:"recommendations"`
}

func CalculateTirePressure(req CalculatorRequest) CalculatorResponse {
	totalWeight := req.RiderWeight + req.BikeWeight

	// Simplified SRAM formula: base pressure from weight and tire width
	// Wider tires = lower pressure, heavier rider = higher pressure
	basePressure := totalWeight / (float64(req.TireWidth) * 1.5)

	// Tire type adjustments
	switch strings.ToLower(req.TireType) {
	case "tubeless":
		basePressure *= 0.93
	case "tubular":
		basePressure *= 0.96
	// clincher is default, no adjustment
	}

	// Surface adjustments
	switch strings.ToLower(req.Surface) {
	case "gravel":
		basePressure *= 0.85
	case "mixed":
		basePressure *= 0.90
	case "cobblestone":
		basePressure *= 0.88
	// road is default, no adjustment
	}

	// Weather/conditions adjustments
	switch strings.ToLower(req.Conditions) {
	case "wet":
		basePressure *= 0.95
	case "mud":
		basePressure *= 0.88
	case "snow":
		basePressure *= 0.85
	}

	// Front/rear split: ~45/55 weight distribution
	frontPressure := math.Round(basePressure*0.95*100) / 100
	rearPressure := math.Round(basePressure*1.05*100) / 100

	// Clamp to reasonable range
	frontPressure = math.Max(1.5, math.Min(frontPressure, 9.0))
	rearPressure = math.Max(1.5, math.Min(rearPressure, 9.0))

	recommendations := buildRecommendations(req, frontPressure, rearPressure)

	return CalculatorResponse{
		FrontPressure:   frontPressure,
		RearPressure:    rearPressure,
		Unit:            "bar",
		Recommendations: recommendations,
	}
}

func buildRecommendations(req CalculatorRequest, front, rear float64) []string {
	var recs []string

	if strings.ToLower(req.TireType) == "tubeless" {
		recs = append(recs, "Tubeless setup allows lower pressures for better grip and comfort.")
	}

	if strings.ToLower(req.Conditions) == "wet" {
		recs = append(recs, "Wet conditions: consider reducing pressure by 0.1-0.2 bar for extra grip.")
	}

	if strings.ToLower(req.Surface) == "gravel" {
		recs = append(recs, "For gravel: lower pressure improves traction. Watch for pinch flats if using tubes.")
	}

	if req.TireWidth >= 35 {
		recs = append(recs, "Wide tires perform best at lower pressures — prioritize comfort and grip.")
	}

	if req.TireWidth <= 25 {
		recs = append(recs, "Narrow tires: be careful not to go too low to avoid pinch flats.")
	}

	recs = append(recs, fmt.Sprintf("Recommended range: front %.1f-%.1f bar, rear %.1f-%.1f bar.",
		math.Max(1.5, front-0.2), front+0.2, math.Max(1.5, rear-0.2), rear+0.2))

	return recs
}
