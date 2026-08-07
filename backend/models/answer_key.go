package models

import "time"

type AnswerKey struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID   string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	OptionID   uint      `json:"option_id" gorm:"not null;index"`
	QuestionID uint      `json:"question_id" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Option   Option   `json:"option,omitempty" gorm:"foreignKey:OptionID"`
	Question Question `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}