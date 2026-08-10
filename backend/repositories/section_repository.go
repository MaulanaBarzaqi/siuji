package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type SectionRepository interface {
	Create(section *models.Section) error
	FindAllPagination(filter, sort string, limit, offset int) ([]models.Section, int64, error)
	FindByPublicID(publicID string) (*models.Section, error)
	Update(section *models.Section) error
	Delete(publicID string) error
	// answer key management
	UpsertAnswerKey(questionPublicID, optionPublicID string) (*models.AnswerKey, error)
	// reorder management
	UpdateQuestionPositions(questionPublicIDs []string) error
	UpdateOptionPositions(optionPublicIDs []string) error
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

func (r *sectionRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.Section, int64, error) {
	var sections []models.Section
	var totalData int64

	query := r.db.Model(&models.Section{})

	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}
	err := query.Limit(limit).Offset(offset).Find(&sections).Error
	return sections, totalData, err
}

func (r *sectionRepository) FindByPublicID(publicID string) (*models.Section, error) {
	var section models.Section
	err := r.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).Preload("Questions.AnswerKeys.Option").Where("public_id = ?", publicID).First(&section).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &section, nil
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
		return tx.Delete(&section).Error
	})
}

func (r *sectionRepository) UpsertAnswerKey(questionPublicID, optionPublicID string) (*models.AnswerKey, error) {
	var answerKey models.AnswerKey
	var question models.Question
	var option models.Option

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Cari Question ID berdasarkan public_id
		if err := tx.Where("public_id = ?", questionPublicID).First(&question).Error; err != nil {
			return errors.New("question not found")
		}
		// Cari Option ID berdasarkan public_id
		if err := tx.Where("public_id = ?", optionPublicID).First(&option).Error; err != nil {
			return errors.New("option not found")
		}
		// Cek apakah answer key sudah ada untuk question ini
		err := tx.Where("question_id = ?", question.ID).First(&answerKey).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				answerKey = models.AnswerKey{
					QuestionID: question.ID,
					OptionID:   option.ID,
				}
				return tx.Create(&answerKey).Error
			}
			return err
		}
		// Update jika sudah ada
		answerKey.OptionID = option.ID
		return tx.Save(&answerKey).Error
	})

	if err != nil {
		return nil, err
	}
	// Load relasi Option agar data yang dikembalikan lengkap
	r.db.Preload("Option").First(&answerKey, answerKey.ID)
	return &answerKey, nil
}

func (r *sectionRepository) UpdateQuestionPositions(questionPublicIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for index, pubID := range questionPublicIDs {
			err := tx.Model(&models.Question{}).
					Where("public_id = ?", pubID).
					Update("position", index+1).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *sectionRepository) UpdateOptionPositions(optionPublicIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for index, pubID := range optionPublicIDs {
			err := tx.Model(&models.Option{}).
				Where("public_id = ?", pubID).
				Update("position", index+1).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}