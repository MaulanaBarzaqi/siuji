package models

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantAnswer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PublicID     uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ExamPeriodID uint      `json:"exam_period_id" gorm:"not null"`
	QuestionID   uint      `json:"question_id" gorm:"not null"`
	OptionID     *uint     `json:"option_id"`
	IsCorrect    bool      `json:"is_correct" gorm:"default:false;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ExamPeriod ExamPeriod `json:"exam_period,omitempty" gorm:"foreignKey:ExamPeriodID"`
	Question   Question   `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Option     *Option    `json:"option,omitempty" gorm:"foreignKey:OptionID"`
}