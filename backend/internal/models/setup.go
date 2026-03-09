package models

import (
	"gorm.io/gorm"
)

type Setup struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id" gorm:"not null;index"`
	Name                  string         `json:"name" gorm:"not null"`
	Photo                 *string        `json:"photo"`
	BikeName              string         `json:"bike_name" gorm:"not null"`
	Tires                 string         `json:"tires" gorm:"not null"`
	ComponentsDescription *string        `json:"components_description" gorm:"type:text"`
	CreatedAt             int64          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             int64          `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}
