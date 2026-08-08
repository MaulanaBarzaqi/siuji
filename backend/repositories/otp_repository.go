package repositories

import (
	"errors"
	"siuji-backend/config"
	"siuji-backend/models"
	"time"

	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(otp *models.OTP) error
	FindValidByEmailAndCode(email, code, purpose string) (*models.OTP, error)
	DeleteByEmail(email string) error
	DeleteExpired() error
}

type otpRepository struct{
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otp *models.OTP) error {
	return config.DB.Create(otp).Error
}

func (r *otpRepository) FindValidByEmailAndCode(email, code, purpose string) (*models.OTP, error) {
	var otp models.OTP
	now := time.Now()

	err := config.DB.
		Where("email = ? AND code = ? AND purpose = ? AND expires_at > ?", email, code, purpose, now).
		First(&otp).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) DeleteByEmail(email string) error {
	return config.DB.
		Where("email = ?", email).
		Delete(&models.OTP{}).Error
}

func (r *otpRepository) DeleteExpired() error {
	now := time.Now()
	return config.DB.
		Where("expires_at < ?", now).
		Delete(&models.OTP{}).Error
}