package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type SectionRepository interface {
	Create(section *models.Section) error
	FindAll() ([]models.Section, error)
	FindByPublicID(publicID string) (*models.Section, error)
	Update(section *models.Section) error
	Delete(publicID string) error
	// answer key management
	UpsertAnswerKey(questionID uint, optionID uint) error
	// reorder management
	UpdateQuestionPositions(updates []QuestionPositionUpdate) error
	UpdateOptionPositions(updates []OptionPositionUpdate) error
}

type QuestionPositionUpdate struct {
	QuestionID uint
	Position   int
}

type OptionPositionUpdate struct {
	OptionID uint
	Position int
}

type sectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{db: db}
}

func (r *sectionRepository) Create(section *models.Section) error {
	return r.db.Create(section).Error
}

func (r *sectionRepository) FindAll() ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Preload("Questions.Options").Find(&sections).Error
	return sections, err
}

func (r *sectionRepository) FindByPublicID(publicID string) (*models.Section, error) {
	var section models.Section
	err := r.db.Preload("Questions.Options").Where("public_id = ?", publicID).First(&section).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &section, err
}

func (r *sectionRepository) Update(section *models.Section) error {
	return r.db.Save(section).Error
}

func (r *sectionRepository) Delete(publicID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var section models.Section
		if err := tx.Where("public_id = ?", publicID).First(&section).Error; err != nil {
			return err
		}
		// Hapus relasi di period_section terlebih dahulu
		if err := tx.Where("section_id = ?", section.ID).Delete(&models.PeriodSection{}).Error; err != nil {
			return err
		}
		// Hapus section
		return tx.Delete(&section).Error
	})
}

func (r *sectionRepository) UpsertAnswerKey(questionID uint, optionID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var answerKey models.AnswerKey
		// Cek apakah kunci jawaban untuk question_id ini sudah ada
		err := tx.Where("question_id = ?", questionID).First(&answerKey).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Jika belum ada, buat baru (generate public_id sesuai kebutuhan/helper uuid Anda)
				newKey := models.AnswerKey{
					QuestionID: questionID,
					OptionID:   optionID,
				}
				return tx.Create(&newKey).Error
			}
			return err
		}
		// Jika sudah ada, update option_id-nya
		answerKey.OptionID = optionID
		return tx.Save(&answerKey).Error
	})
}

func (r *sectionRepository) UpdateQuestionPositions(updates []QuestionPositionUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			err := tx.Model(&models.Question{}).
				Where("id = ?", update.QuestionID).
				Update("position", update.Position).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *sectionRepository) UpdateOptionPositions(updates []OptionPositionUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			err := tx.Model(&models.Option{}).
				Where("id = ?", update.OptionID).
				Update("position", update.Position).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}