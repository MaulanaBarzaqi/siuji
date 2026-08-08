package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(question *models.Question) error
	FindByPublicID(publicID string) (*models.Question, error)
	Update(question *models.Question) error
	Delete(publicID string) error
	GetMaxPositionInSection(sectionID uint) (int, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) FindByPublicID(publicID string) (*models.Question, error) {
	var question models.Question
	err := r.db.Preload("Options").Where("public_id = ?", publicID).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &question, err
}

func (r *questionRepository) Update(question *models.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) Delete(publicID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var question models.Question
		if err := tx.Where("public_id = ?", publicID).First(&question).Error; err != nil {
			return err
		}
		// hapus kunci jawaban soal ini
		if err := tx.Where("question_id = ?", question.ID).Delete(&models.AnswerKey{}).Error; err != nil {
			return err
		}
		// hapus semua opsi soal ini
		if err := tx.Where("question_id = ?", question.ID).Delete(&models.Option{}).Error; err != nil {
			return err
		}
		return tx.Delete(&question).Error
	})
}

func (r *questionRepository) GetMaxPositionInSection(sectionID uint) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.Question{}).
			Where("section_id = ?", sectionID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&maxPosition).Error
	return maxPosition, err
}
