package models

import "time"

type Option struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID   string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	QuestionID uint      `json:"question_id" gorm:"not null;index"`
	Label      string    `json:"label" gorm:"type:varchar(10);not null"`
	OptionText string    `json:"option_text" gorm:"type:text;not null"`
	Position   int       `json:"position" gorm:"type:int;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Question           Question            `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	AnswerKeys         []AnswerKey         `json:"answer_keys,omitempty" gorm:"foreignKey:OptionID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:OptionID"`
}