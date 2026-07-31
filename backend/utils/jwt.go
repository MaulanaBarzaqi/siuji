package utils

import (
	"errors"
	"fmt"
	"siuji-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	PurposeSetupCredential = "setup_credential"
	PurposeVerifyEmail = "verify_email"
	PurposeResetPassword = "reset_password"
)

type AccessTokenClaims struct {
	UserID		uint		`json:"user_id"`
	PublicID	uuid.UUID	`json:"pub_id"`
	Email		string		`json:"email"`
	Role		string		`json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID		uint	`json:"user_id"`
	TokenFamily	string	`json:"token_family"`
	jwt.RegisteredClaims
}

type TempTokenClaims struct {
	UserID		uint	`json:"user_id"`
	Purpose		string	`json:"purpose"`
	Email		string	`json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint, publicID uuid.UUID, email, role string) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET required")
	}
	duration, err := time.ParseDuration(config.AppConfig.JWTExpiresIn)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRES_IN format: %w", err)
	}
	now := time.Now()
	claims := AccessTokenClaims{
		UserID: userID,
		PublicID: publicID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer: "siuji-backend",
			Subject: fmt.Sprintf("%d", userID),
			ID: uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint, tokenFamily string) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET required")
	}
	duration, err := time.ParseDuration(config.AppConfig.RefreshTokenExpires)
	if err != nil {
		return "", fmt.Errorf("invalid REFRESH_TOKEN_EXPIRES format: %w", err)
	}
	now := time.Now()
	claims := RefreshTokenClaims{
		UserID: userID,
		TokenFamily: tokenFamily,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer: "siuji-backend",
			Subject: fmt.Sprintf("%d", userID),
			ID: uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateTempToken(userID uint, email, purpose string, expiryMinutes int) (string, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return "", errors.New("JWT_SECRET required")
	}
	validPurposes := map[string]bool{
		PurposeSetupCredential: true,
		PurposeVerifyEmail: true,
		PurposeResetPassword: true,
	}
	if !validPurposes[purpose] {
		return "", fmt.Errorf("invalid purpose: %s", purpose)
	}
	now := time.Now()
	claims := TempTokenClaims{
		UserID: userID,
		Email: email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer: "siuji-backend",
			Subject: fmt.Sprintf("%d", userID),
			ID: uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	secret := config.AppConfig.JWTSecret

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token)(interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	secret := config.AppConfig.JWTSecret

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token)(interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func ValidateTempToken(tokenString string, expectedPurpose string) (*TempTokenClaims, error) {
	secret := config.AppConfig.JWTSecret
	token, err := jwt.ParseWithClaims(tokenString, &TempTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*TempTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	if claims.Purpose != expectedPurpose {
		return nil, fmt.Errorf("invalid token purpose: expected %s, got %s", expectedPurpose, claims.Purpose)
	}
	return claims, nil
}

func GenerateTokenFamily() string {
	return uuid.New().String()
}

