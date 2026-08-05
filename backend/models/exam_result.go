package models

import (
	"time"

	"github.com/google/uuid"
)

type ExamResult struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	PublicID         uuid.UUID  `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	UserID           uint       `json:"user_id" gorm:"not null"`
	ExamPeriodID     uint       `json:"exam_period_id" gorm:"not null"`
	TotalScaledScore float64    `json:"total_scaled_score" gorm:"default:0;not null"`
	Status           string     `json:"status" gorm:"type:varchar(50);not null"` // in_progress, completed, passed, failed
	SubmittedAt      *time.Time `json:"submitted_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relasi
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ExamPeriod ExamPeriod `json:"exam_period,omitempty" gorm:"foreignKey:ExamPeriodID"`
}