package services

import (
	"github.com/levferril/racenotes/backend/internal/models"
	"gorm.io/gorm"
)

type RaceService struct {
	DB *gorm.DB
}

func NewRaceService(db *gorm.DB) *RaceService {
	return &RaceService{DB: db}
}

func (s *RaceService) Create(race *models.Race) error {
	if race.SetupID != nil {
		var setup models.Setup
		if err := s.DB.Where("id = ? AND user_id = ?", *race.SetupID, race.UserID).First(&setup).Error; err != nil {
			return err
		}
		race.BikeName = &setup.BikeName
		race.Tires = &setup.Tires
	}
	return s.DB.Create(race).Error
}

func (s *RaceService) List(userID uint, raceType, isCompleted, setupID string) ([]models.Race, error) {
	var races []models.Race
	query := s.DB.Where("user_id = ?", userID).Preload("Setup")

	if raceType != "" {
		query = query.Where("type = ?", raceType)
	}
	if isCompleted != "" {
		query = query.Where("is_completed = ?", isCompleted == "true")
	}
	if setupID != "" {
		query = query.Where("setup_id = ?", setupID)
	}

	err := query.Order("date desc").Find(&races).Error
	return races, err
}

func (s *RaceService) GetByID(userID, raceID uint) (*models.Race, error) {
	var race models.Race
	err := s.DB.Where("id = ? AND user_id = ?", raceID, userID).Preload("Setup").First(&race).Error
	if err != nil {
		return nil, err
	}
	return &race, nil
}

func (s *RaceService) Update(userID, raceID uint, updates map[string]interface{}) (*models.Race, error) {
	var race models.Race
	if err := s.DB.Where("id = ? AND user_id = ?", raceID, userID).First(&race).Error; err != nil {
		return nil, err
	}

	if sid, ok := updates["setup_id"]; ok && sid != nil {
		setupID := uint(sid.(float64))
		var setup models.Setup
		if err := s.DB.Where("id = ? AND user_id = ?", setupID, userID).First(&setup).Error; err != nil {
			return nil, err
		}
		updates["bike_name"] = setup.BikeName
		updates["tires"] = setup.Tires
	}

	if err := s.DB.Model(&race).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetByID(userID, raceID)
}

func (s *RaceService) Delete(userID, raceID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", raceID, userID).Delete(&models.Race{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
