package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID	     		uint 			`json:"id" gorm:"primaryKey"`
	PublicId 		uuid.UUID		`json:"public_id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name			string			`json:"name" gorm:"not null"`
	Email			string			`json:"email" gorm:"unique;not null"`
	University		string			`json:"university" gorm:"default:null"`
	NIM				string 			`json:"nim" gorm:"uniqueIndex;default:null"`
	Password		string			`json:"-" gorm:"column:password;not null"`
	Role			string			`json:"role" gorm:"default:participant;not null"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"-" gorm:"index"`
}