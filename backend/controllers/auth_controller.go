package controllers

import (
	"errors"
	"siuji-backend/services"
	"siuji-backend/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{service: s}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyOTPRequest struct {
	OTPCode string `json:"otp_code"`
}

type ResetPasswordRequest struct {
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func getUserIDFromContext(ctx fiber.Ctx) (uint, error) {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok {
		return 0, errors.New("invalid user context")
	}
	return userID, nil
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} map[string]interface{} "Successful login"
// @Failure      400 {object} map[string]interface{} "Invalid request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login(ctx fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	response, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login failed", err.Error())
	}

	if response.AccessToken != "" {
		utils.SetAuthCookies(ctx, response.AccessToken, response.RefreshToken)
	}

	return utils.Success(ctx, "Login successfully. Welcome to the dashboard.", response)
}

// ForgotPassword godoc
// @Summary      Request password reset OTP
// @Description  Send an OTP verification code to the registered email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body ForgotPasswordRequest true "User email"
// @Success      200 {object} map[string]interface{} "OTP sent successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /api/v1/auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	response, err := c.service.ForgotPassword(req.Email)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to process forgot password", err.Error())
	}

	if response == nil {
		return utils.SuccessNoData(ctx, "If this email is registered, an OTP verification code has been sent.")
	}

	return utils.Success(ctx, "OTP verification code has been sent to your email.", response)
}

// VerifyOTP godoc
// @Summary      Verify OTP code
// @Description  Verify the OTP code sent to email for password reset
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body VerifyOTPRequest true "OTP code"
// @Success      200 {object} map[string]interface{} "OTP verified successfully"
// @Failure      400 {object} map[string]interface{} "Verification failed"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /api/v1/auth/verify-otp [post]
func (c *AuthController) VerifyOTP(ctx fiber.Ctx) error {
	var req VerifyOTPRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized", err.Error())
	}

	response, err := c.service.VerifyOTP(userID, req.OTPCode)
	if err != nil {
		return utils.BadRequest(ctx, "OTP verification failed", err.Error())
	}

	return utils.Success(ctx, "OTP verified. Please set your new password.", response)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using new password after OTP verification
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body ResetPasswordRequest true "New password details"
// @Success      200 {object} map[string]interface{} "Password reset successful"
// @Failure      400 {object} map[string]interface{} "Reset failed"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /api/v1/auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized", err.Error())
	}

	err = c.service.ResetPassword(userID, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		return utils.BadRequest(ctx, "Password reset failed", err.Error())
	}

	return utils.SuccessNoData(ctx, "Password reset successfully. Please login with your new password.")
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change current user password by providing old and new password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body ChangePasswordRequest true "Old and new password"
// @Success      200 {object} map[string]interface{} "Password changed successfully"
// @Failure      400 {object} map[string]interface{} "Change failed"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /api/v1/auth/change-password [post]
func (c *AuthController) ChangePassword(ctx fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized", err.Error())
	}

	err = c.service.ChangePassword(userID, req.OldPassword, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		return utils.BadRequest(ctx, "Password change failed", err.Error())
	}

	return utils.Success(ctx, "Password changed successfully.", nil)
}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Retrieve logged-in user profile details
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]interface{} "User profile retrieved"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /api/v1/auth/me [get]
func (c *AuthController) GetMe(ctx fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized", err.Error())
	}

	response, err := c.service.GetMe(userID)
	if err != nil {
		return utils.NotFound(ctx, "User not found", err.Error())
	}

	return utils.Success(ctx, "User profile retrieved successfully.", response)
}

// Logout godoc
// @Summary      User logout
// @Description  Clear authentication cookies and log out user
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]interface{} "Logout successful"
// @Router       /api/v1/auth/logout [post]
func (c *AuthController) Logout(ctx fiber.Ctx) error {
	utils.ClearAuthCookies(ctx)
	return utils.Success(ctx, "Logout successfully.", nil)
}