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
	Password		string			`json:"-" gorm:"column:password;not null"`
	Role			string			`json:"role" gorm:"default:user;not null"`
	IsFirstLogin	bool			`json:"is_first_login" gorm:"default:true;not null"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"-" gorm:"index"`
}