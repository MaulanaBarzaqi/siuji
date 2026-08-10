package models

import "time"

type Section struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	PeriodSections []PeriodSection `json:"period_sections,omitempty" gorm:"foreignKey:SectionID"`
	Questions      []Question      `json:"questions,omitempty" gorm:"foreignKey:SectionID"`
	SectionScores  []SectionScore  `json:"section_scores,omitempty" gorm:"foreignKey:SectionID"`
}

type SectionRequest struct {
	Title string `json:"title" validate:"required"`
}

type SectionDetailResponse struct {
	PublicID  string                   `json:"public_id"`
	Title     string                   `json:"title"`
	Questions []QuestionDetailResponse `json:"questions"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}
