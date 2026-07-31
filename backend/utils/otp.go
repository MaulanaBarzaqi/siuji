package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"siuji-backend/models"
	"time"
)

const (
	OTPPurposeResetPassword = "reset-password"
)

func GenerateOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func NewOTP(email, purpose string) (*models.OTP, error) {
	code, err := GenerateOTPCode()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	otp := &models.OTP{
		Email: email,
		Code: code,
		Purpose: purpose,
		ExpiresAt: now.Add(5 *time.Minute),
	}
	return otp, nil
}