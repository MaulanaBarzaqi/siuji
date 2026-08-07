package models

import "time"

type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	SectionID uint      `json:"section_id" gorm:"not null;index"`
	Question  string    `json:"question" gorm:"type:text;not null"`
	AudioURL  *string   `json:"audio_url" gorm:"type:text"`
	ImageURL  *string   `json:"image_url" gorm:"type:text"`
	Passage   *string   `json:"passage" gorm:"type:text"`
	Position  int       `json:"position" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Section            Section             `json:"section,omitempty" gorm:"foreignKey:SectionID"`
	Options            []Option            `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	AnswerKeys         []AnswerKey         `json:"answer_keys,omitempty" gorm:"foreignKey:QuestionID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:QuestionID"`
}