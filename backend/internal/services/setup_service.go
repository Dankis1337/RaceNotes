package services

import (
	"github.com/levferril/racenotes/backend/internal/models"
	"gorm.io/gorm"
)

type SetupService struct {
	DB *gorm.DB
}

func NewSetupService(db *gorm.DB) *SetupService {
	return &SetupService{DB: db}
}

func (s *SetupService) Create(setup *models.Setup) error {
	return s.DB.Create(setup).Error
}

func (s *SetupService) List(userID uint) ([]models.Setup, error) {
	var setups []models.Setup
	err := s.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&setups).Error
	return setups, err
}

func (s *SetupService) GetByID(userID, setupID uint) (*models.Setup, error) {
	var setup models.Setup
	err := s.DB.Where("id = ? AND user_id = ?", setupID, userID).First(&setup).Error
	if err != nil {
		return nil, err
	}
	return &setup, nil
}

func (s *SetupService) Update(userID, setupID uint, updates map[string]interface{}) (*models.Setup, error) {
	var setup models.Setup
	if err := s.DB.Where("id = ? AND user_id = ?", setupID, userID).First(&setup).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Model(&setup).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &setup, nil
}

func (s *SetupService) Delete(userID, setupID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", setupID, userID).Delete(&models.Setup{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
