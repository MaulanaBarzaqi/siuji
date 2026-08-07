package models

type ScoreConversion struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	SectionType  string `json:"section_type" gorm:"type:varchar(50);not null;index"`
	CorrectCount int    `json:"correct_count" gorm:"type:int;not null"`
	ScaledScore  int    `json:"scaled_score" gorm:"type:int;not null"`
}