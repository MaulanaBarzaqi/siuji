package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type OptionRepository interface {
	Create(option *models.Option) error
	FindByPublicID(publicID string) (*models.Option, error)
	Update(option *models.Option) error
	Delete(publicID string) error
	GetMaxPositionInQuestion(questionID uint) (int, error)
}

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) Create(option *models.Option) error {
	return r.db.Create(option).Error
}

func (r *optionRepository) FindByPublicID(publicID string) (*models.Option, error) {
	var option models.Option
	err := r.db.Where("public_id = ?", publicID).First(&option).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &option, err
}

func (r *optionRepository) Update(option *models.Option) error {
	return r.db.Save(option).Error
} 

func (r *optionRepository) Delete(publicID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var option models.Option
		if err := tx.Where("public_id = ?", publicID).First(&option).Error; err != nil {
			return err
		}
		// jika opsi terdaftar sbg kunci jawaban, hapus juga record answer key-nya
		if err := tx.Where("option_id = ?", option.ID).Delete(&models.AnswerKey{}).Error; err != nil {
			return err
		}
		return tx.Delete(&option).Error
	})
	
}

func (r *optionRepository) GetMaxPositionInQuestion(questionID uint) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.Option{}).
			Where("question_id = ?", questionID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&maxPosition).Error
	return maxPosition, err
}
