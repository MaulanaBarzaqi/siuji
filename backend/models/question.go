package models

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PublicID  uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid();not null"`
	SectionID uint      `json:"section_id" gorm:"not null"`
	Question  string    `json:"question" gorm:"type:text;not null"`
	AudioURL  string    `json:"audio_url" gorm:"type:text"`
	ImageURL  string    `json:"image_url" gorm:"type:text"`
	Passage   string    `json:"passage" gorm:"type:text"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi
	Section            ExamSection         `json:"section,omitempty" gorm:"foreignKey:SectionID"`
	Options            []Option            `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:QuestionID"`
}