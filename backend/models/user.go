package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID	     		uint 			`json:"id" gorm:"primaryKey"`
	PublicID 		uuid.UUID		`json:"public_id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name			string			`json:"name" gorm:"not null"`
	Email			string			`json:"email" gorm:"unique;not null"`
	University		string			`json:"university" gorm:"default:null"`
	NIM				string 			`json:"nim" gorm:"uniqueIndex;default:null"`
	Password		string			`json:"-" gorm:"column:password;not null"`
	Role			string			`json:"role" gorm:"default:participant;not null"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"-" gorm:"index"`

	// Relations
	ParticipantPeriods []ParticipantPeriod `json:"participant_periods,omitempty" gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	PublicID   string    `json:"public_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	NIM        string    `json:"nim,omitempty"`
	University string    `json:"university,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type AuthResponse struct {
	TempToken    string        `json:"temp_token,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresIn    int           `json:"expires_in,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		PublicID:   u.PublicID.String(),
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		NIM:        u.NIM,
		University: u.University,
		UpdatedAt:  u.UpdatedAt,
	}
}