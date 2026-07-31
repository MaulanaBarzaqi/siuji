package services

import "log"

type EmailService interface {
	SendOTP(email, otpCode string) error
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) SendOTP(email, otpCode string) error {
	log.Printf("OTP EMAIL SIMULATION")
    log.Printf("To: %s", email)
    log.Printf("OTP Code: %s", otpCode)
    return nil
}