package models

type ScoreConversion struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	SectionID   uint `json:"section_id" gorm:"not null"`
	RawCount    int  `json:"raw_count" gorm:"not null"`
	ScaledScore int  `json:"scaled_score" gorm:"not null"`

	// Relasi
	Section ExamSection `json:"section,omitempty" gorm:"foreignKey:SectionID"`
}