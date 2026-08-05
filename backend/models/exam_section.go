package models

import (
	"time"

	"github.com/google/uuid"
)

type ExamSection struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PublicID     uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	ExamPeriodID uint      `json:"exam_period_id" gorm:"not null"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi
	ExamPeriod         ExamPeriod          `json:"exam_period,omitempty" gorm:"foreignKey:ExamPeriodID"`
	Questions          []Question          `json:"questions,omitempty" gorm:"foreignKey:SectionID"`
	ScoreConversions   []ScoreConversion   `json:"score_conversions,omitempty" gorm:"foreignKey:SectionID"`
	ExamResultSections []ExamResultSection `json:"exam_result_sections,omitempty" gorm:"foreignKey:SectionID"`
}