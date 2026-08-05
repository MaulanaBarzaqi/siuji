package models

import (
	"time"

	"github.com/google/uuid"
)

type ExamResultSection struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	PublicID       uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	ExamPeriodID   uint      `json:"exam_period_id" gorm:"not null"`
	SectionID      uint      `json:"section_id" gorm:"not null"`
	CorrectCount   int       `json:"correct_count" gorm:"default:0;not null"`
	ConvertedScore int       `json:"converted_score" gorm:"default:0;not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relasi
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ExamPeriod ExamPeriod  `json:"exam_period,omitempty" gorm:"foreignKey:ExamPeriodID"`
	Section    ExamSection `json:"section,omitempty" gorm:"foreignKey:SectionID"`
}