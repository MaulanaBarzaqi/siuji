package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamPeriod struct {
	ID				       uint 		  `json:"id" gorm:"primaryKey"`
	PublicID 			   uuid.UUID	  `json:"public_id" gorm:"type:uuid;default:gen_random_uuid()"`
	Title	 			   string		  `json:"title" gorm:"type:varchar(255);not null"`
	Month                  int            `json:"month" gorm:"not null"`
	Year                   int            `json:"year" gorm:"not null"`
	Status                 string         `json:"status" gorm:"type:varchar(50);not null"`
	CertificateExpMonths   int            `json:"certificate_exp_months"`
	CertificateTemplateURL string         `json:"certificate_template_url" gorm:"type:text"`
	PassingGrade           float64        `json:"passing_grade"`
	StartTime              time.Time      `json:"start_time"`
	EndTime                time.Time      `json:"end_time"`
	CreatedAt              time.Time      `json:"created_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`

	// relasi
	Sections               []ExamSection       `json:"sections,omitempty" gorm:"foreignKey:ExamPeriodID"`
	ParticipantAnswers     []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:ExamPeriodID"`
	ExamResults            []ExamResult        `json:"exam_results,omitempty" gorm:"foreignKey:ExamPeriodID"`
	ExamResultSections     []ExamResultSection `json:"exam_result_sections,omitempty" gorm:"foreignKey:ExamPeriodID"`
}