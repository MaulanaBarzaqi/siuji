package models

import "time"

type Period struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID            string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	Title               string    `json:"title" gorm:"type:varchar(255);not null"`
	Month               string    `json:"month" gorm:"type:varchar(50)"`
	Year                int       `json:"year" gorm:"type:int"`
	Status              string    `json:"status" gorm:"type:varchar(50);not null"`
	CertificateURL      string    `json:"certificate_url" gorm:"type:text"`
	CertificateExpMonth int       `json:"certificate_exp_month" gorm:"type:int"`
	MinPassingGrade     int       `json:"min_passing_grade" gorm:"type:int"`
	MaxPassingGrade     int       `json:"max_passing_grade" gorm:"type:int"`
	StartTime           time.Time `json:"start_time" gorm:"type:timestamp"`
	EndTime             time.Time `json:"end_time" gorm:"type:timestamp"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	ParticipantPeriods []ParticipantPeriod `json:"participant_periods,omitempty" gorm:"foreignKey:PeriodID"`
	PeriodSections     []PeriodSection     `json:"period_sections,omitempty" gorm:"foreignKey:PeriodID"`
}