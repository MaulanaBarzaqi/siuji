package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type ScoreConversionRepository interface {
	Create(scoreConversion *models.ScoreConversion) error
	FindAll() ([]models.ScoreConversion, error)
	FindBySectionTypeAndCorrectCount(sectionType string, correctCount int) (*models.ScoreConversion, error)
	Update(scoreConversion *models.ScoreConversion) error
	Delete(id uint) error
}

type scoreConversionRepository struct {
	db *gorm.DB
}

func NewScoreConversionRepository(db *gorm.DB) ScoreConversionRepository {
	return &scoreConversionRepository{db: db}
}

func (r *scoreConversionRepository) Create(scoreConversion *models.ScoreConversion) error {
	return r.db.Create(scoreConversion).Error
}

func (r *scoreConversionRepository) FindAll() ([]models.ScoreConversion, error) {
	var conversions []models.ScoreConversion
	err := r.db.Find(&conversions).Error
	return conversions, err
}

func (r *scoreConversionRepository) FindBySectionTypeAndCorrectCount(sectionType string, correctCount int) (*models.ScoreConversion, error) {
	var conversion models.ScoreConversion
	err := r.db.Where("section_type = ? AND correct_count = ?", sectionType, correctCount).First(&conversion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conversion, nil
}

func (r *scoreConversionRepository) Update(scoreConversion *models.ScoreConversion) error {
	return r.db.Save(scoreConversion).Error
}

func (r *scoreConversionRepository) Delete(id uint) error {
	return r.db.Delete(&models.ScoreConversion{}, id).Error
}