package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"
	"siuji-backend/utils"
	"time"
)

const (
	TempTokenExpiryMinutes = 5
	AccessTokenExpiryHours = 1
)

type AuthService interface {
	Login(email, password string) (*AuthResponse, error)
	SetupCredential(userID uint, newName, newEmail, newPassword, confirmNewPassword string) (*AuthResponse, error)
	ForgotPassword(email string) (*AuthResponse, error)
	VerifyOTP(userID uint, otpCode string) (*AuthResponse, error)
	ResetPassword(userID uint, newPassword, confirmNewPassword string) error
	ChangePassword(userID uint, oldPassword, newPassword, confirmNewPassword string) error
	GetMe(userID uint) (*UserResponse, error)
}

type authService struct {
	userRepo     repositories.UserRepository
	otpRepo      repositories.OTPRepository
	emailService EmailService
}

func NewAuthService(
	userRepo repositories.UserRepository,
	otpRepo repositories.OTPRepository,
	emailService EmailService,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		emailService: emailService,
	}
}

type UserResponse struct {
	PublicID     string    `json:"public_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsFirstLogin *bool     `json:"is_first_login,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type AuthResponse struct {
	TempToken    string        `json:"temp_token,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresIn    int           `json:"expires_in,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
}

func buildUserResponse(user *models.User, includeFirstLogin bool) *UserResponse {
	res := &UserResponse{
		PublicID:  user.PublicId.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		UpdatedAt: user.UpdatedAt,
	}

	if includeFirstLogin {
		res.IsFirstLogin = &user.IsFirstLogin
	}
	return res
}


func (s *authService) Login(email, password string) (*AuthResponse, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	if user.IsFirstLogin {
		tempToken, err := utils.GenerateTempToken(user.ID, user.Email, utils.PurposeSetupCredential, TempTokenExpiryMinutes)
		if err != nil {
			return nil, errors.New("failed to generate temp token")
		}
		return &AuthResponse{
			TempToken: tempToken,
			ExpiresIn: TempTokenExpiryMinutes * 60,
			User:      buildUserResponse(user, true),
		}, nil
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.PublicId, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, utils.GenerateTokenFamily())
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    AccessTokenExpiryHours * 3600,
		User:         buildUserResponse(user, false),
	}, nil
}

func (s *authService) SetupCredential(userID uint, newName, newEmail, newPassword, confirmNewPassword string) (*AuthResponse, error) {
	if newName == "" || newEmail == "" {
		return nil, errors.New("name and email are required")
	}
	if err := utils.ValidatePassword(newPassword, confirmNewPassword); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !user.IsFirstLogin {
		return nil, errors.New("credential already setup")
	}
	if newEmail != user.Email && s.userRepo.EmailExists(newEmail) {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	if err := s.userRepo.UpdateProfile(userID, newName, newEmail); err != nil {
		return nil, errors.New("failed to update profile")
	}
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return nil, errors.New("failed to update password")
	}
	if err := s.userRepo.UpdateFirstLogin(userID); err != nil {
		return nil, errors.New("failed to update first login status")
	}

	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("failed to get updated user")
	}

	return &AuthResponse{
		User: buildUserResponse(updatedUser, true),
	}, nil
}

func (s *authService) ForgotPassword(email string) (*AuthResponse, error) {
	// SECURITY: Blind Response (Mencegah Account Enumeration)
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil
	}
	if err := s.otpRepo.DeleteByEmail(email); err != nil {
		return nil, errors.New("failed to cleanup OTP")
	}

	otp, err := utils.NewOTP(email, utils.OTPPurposeResetPassword)
	if err != nil {
		return nil, errors.New("failed to generate OTP")
	}
	if err := s.otpRepo.Create(otp); err != nil {
		return nil, errors.New("failed to save OTP")
	}
	if err := s.emailService.SendOTP(email, otp.Code); err != nil {
		return nil, errors.New("failed to send OTP email")
	}

	tempToken, err := utils.GenerateTempToken(user.ID, email, utils.PurposeVerifyEmail, TempTokenExpiryMinutes)
	if err != nil {
		return nil, errors.New("failed to generate temp token")
	}

	return &AuthResponse{
		TempToken: tempToken,
		ExpiresIn: TempTokenExpiryMinutes * 60,
	}, nil
}

func (s *authService) VerifyOTP(userID uint, otpCode string) (*AuthResponse, error) {
	if otpCode == "" {
		return nil, errors.New("OTP code is required")
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	_, err = s.otpRepo.FindValidByEmailAndCode(user.Email, otpCode, utils.OTPPurposeResetPassword)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}
	if err := s.otpRepo.DeleteByEmail(user.Email); err != nil {
		return nil, errors.New("failed to cleanup OTP")
	}
	newTempToken, err := utils.GenerateTempToken(user.ID, user.Email, utils.PurposeResetPassword, TempTokenExpiryMinutes)
	if err != nil {
		return nil, errors.New("failed to generate temp token")
	}
	return &AuthResponse{
		TempToken: newTempToken,
		ExpiresIn: TempTokenExpiryMinutes * 60,
	}, nil
}

func (s *authService) ResetPassword(userID uint, newPassword, confirmNewPassword string) error {
	if err := utils.ValidatePassword(newPassword, confirmNewPassword); err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return errors.New("failed to update password")
	}
	_ = s.otpRepo.DeleteByEmail(user.Email)

	return nil
}

func (s *authService) ChangePassword(userID uint, oldPassword, newPassword, confirmNewPassword string) error {
	if err := utils.ValidatePassword(newPassword, confirmNewPassword); err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

func (s *authService) GetMe(userID uint) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return buildUserResponse(user, false), nil
}