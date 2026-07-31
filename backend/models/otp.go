package models

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	ID        uint   			`json:"id" gorm:"primaryKey"`
	Email     string 			`json:"email" gorm:"index;not null"`
	Code      string 			`json:"code" gorm:"not null"`
	Purpose   string 			`json:"purpose" gorm:"not null"`
	ExpiresAt time.Time 		`json:"expires_at" gorm:"not null"`
	CreatedAt time.Time 		`json:"created_at"`
	UpdatedAt time.Time 		`json:"updated_at"`
	DeletedAt gorm.DeletedAt 	`json:"-" gorm:"index"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}