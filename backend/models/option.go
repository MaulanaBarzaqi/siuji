package models

import (
	"time"

	"github.com/google/uuid"
)

type Option struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PublicID   uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	Label      string    `json:"label" gorm:"type:varchar(10);not null"`
	OptionText string    `json:"option_text" gorm:"type:text;not null"`
	IsCorrect  bool      `json:"is_correct" gorm:"default:false;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relasi
	Question           Question            `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:OptionID"`
}