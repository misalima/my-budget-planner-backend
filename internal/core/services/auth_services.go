package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/misalima/my-budget-planner-backend/internal/api/http/auth"
	"github.com/misalima/my-budget-planner-backend/internal/core/domain"
	"github.com/misalima/my-budget-planner-backend/internal/core/interfaces/irepository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	authRepo irepository.AuthLoader
	userRepo irepository.UserLoader
}

func NewAuthService(authRepo irepository.AuthLoader, userRepo irepository.UserLoader) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}

	err = s.CheckPasswords(password, user.Password)
	if err != nil {
		return "", "", err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	err = s.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *AuthService) SaveRefreshToken(userId uuid.UUID, refreshToken string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.authRepo.StoreRefreshToken(ctx, userId, refreshToken)
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, userId uuid.UUID, token string) (domain.RefreshToken, error) {

	refreshToken, err := s.authRepo.GetRefreshToken(ctx, token)
	if err != nil {
		return refreshToken, fmt.Errorf("invalid refresh token")
	}

	if refreshToken.UserID != userId {
		return refreshToken, fmt.Errorf("refresh token doesnt belong to requesting user")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return refreshToken, fmt.Errorf("refresh token has expired")
	}

	return refreshToken, nil
}

func (s *AuthService) DeleteRefreshToken(token string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.authRepo.DeleteRefreshToken(ctx, token)
}

func (s *AuthService) RefreshToken(userId uuid.UUID, token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	refreshToken, err := s.ValidateRefreshToken(ctx, userId, token)
	if err != nil {
		return "", err
	}

	newAccessToken, err := auth.GenerateAccessToken(refreshToken.UserID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) CheckPasswords(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid email or password")
	}
	return nil
}
