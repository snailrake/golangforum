package usecase

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golangforum/internal/domain"
	"golangforum/internal/repository/postgres"
	"golangforum/internal/utils"
	"time"
)

type AuthUseCase struct {
	Repo *postgres.Repository
}

func NewAuthUseCase(repo *postgres.Repository) *AuthUseCase {
	return &AuthUseCase{
		Repo: repo,
	}
}

func (uc *AuthUseCase) Register(user *domain.User) error {
	existingUser, err := uc.Repo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Role = "user"
	return uc.Repo.CreateUser(user)
}

func (uc *AuthUseCase) Login(username, password string) (string, string, error) {
	user, err := uc.Repo.GetUserByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}
	refreshToken, refreshExp, err := utils.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	if err = uc.Repo.DeleteRefreshTokensByUserID(user.ID); err != nil {
		return "", "", err
	}

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExp,
	}
	if err = uc.Repo.SaveRefreshToken(rt); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (uc *AuthUseCase) RefreshToken(oldToken string) (string, string, error) {
	claims, err := utils.VerifyToken(oldToken)
	if err != nil {
		return "", "", err
	}

	rt, err := uc.Repo.GetRefreshToken(oldToken)
	if err != nil || rt.UserID != claims.UserID || time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("invalid refresh token")
	}

	if err = uc.Repo.DeleteRefreshToken(oldToken); err != nil {
		return "", "", err
	}

	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, newRefreshExp, err := utils.GenerateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}
	if err = uc.Repo.SaveRefreshToken(&domain.RefreshToken{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: newRefreshExp,
	}); err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}
