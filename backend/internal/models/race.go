package models

import (
	"gorm.io/gorm"
)

type Race struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UserID            uint           `json:"user_id" gorm:"not null;index"`
	Name              string         `json:"name" gorm:"not null"`
	Date              string         `json:"date" gorm:"not null"`
	Type              string         `json:"type" gorm:"not null"`
	Photo             *string        `json:"photo"`
	SetupID           *uint          `json:"setup_id"`
	BikeName          *string        `json:"bike_name"`
	Tires             *string        `json:"tires"`
	TirePressureFront *float64       `json:"tire_pressure_front"`
	TirePressureRear  *float64       `json:"tire_pressure_rear"`
	OtherComponents   *string        `json:"other_components" gorm:"type:text"`
	Temperature       *int           `json:"temperature"`
	Conditions        *string        `json:"conditions"`
	Wind              *string        `json:"wind"`
	RoadConditions    *string        `json:"road_conditions"`
	NutritionPlan     *string        `json:"nutrition_plan" gorm:"type:text"`
	Result            *string        `json:"result"`
	Rating            *int           `json:"rating"`
	Feelings          *string        `json:"feelings" gorm:"type:text"`
	IsCompleted       bool           `json:"is_completed" gorm:"default:false"`
	Setup             *Setup         `json:"setup,omitempty" gorm:"foreignKey:SetupID"`
	CreatedAt         int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         int64          `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}
