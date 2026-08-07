package models

import "time"

type PeriodSection struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	PeriodID  uint      `json:"period_id" gorm:"not null;index"`
	SectionID uint      `json:"section_id" gorm:"not null;index"`
	Position  int       `json:"position" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Period  Period  `json:"period,omitempty" gorm:"foreignKey:PeriodID"`
	Section Section `json:"section,omitempty" gorm:"foreignKey:SectionID"`
}