package models

import "time"

type ParticipantPeriod struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	PeriodID  uint      `json:"period_id" gorm:"not null;index"`
	Status    string    `json:"status" gorm:"type:varchar(50);not null"`
	Score     *int      `json:"score" gorm:"type:int"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	User               User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Period             Period              `json:"period,omitempty" gorm:"foreignKey:PeriodID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
	SectionScores      []SectionScore      `json:"section_scores,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
}