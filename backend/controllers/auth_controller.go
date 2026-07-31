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

type SetupCredentialRequest struct {
	NewName               string `json:"new_name"`
	NewEmail           string `json:"new_email"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
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

	message := "Login successfully. Welcome to the dashboard."
	if response.User != nil && response.User.IsFirstLogin != nil && *response.User.IsFirstLogin {
		message = "Login successfully, please setup your credential account."
	}

	return utils.Success(ctx, message, response)
}

func (c *AuthController) SetupCredential(ctx fiber.Ctx) error {
	var req SetupCredentialRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request", err.Error())
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return utils.Unauthorized(ctx, "Unauthorized", err.Error())
	}

	response, err := c.service.SetupCredential(userID, req.NewName, req.NewEmail, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		return utils.BadRequest(ctx, "Setup credential failed", err.Error())
	}

	return utils.Success(ctx, "Credential account setup successfully, please relogin.", response)
}

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

func (c *AuthController) Logout(ctx fiber.Ctx) error {
	utils.ClearAuthCookies(ctx)
	return utils.Success(ctx, "Logout successfully.", nil)
}