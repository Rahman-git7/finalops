package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/finalops/auth-service/internal/model"
	"github.com/finalops/auth-service/internal/repository"
	"github.com/finalops/auth-service/pkg/jwtutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *repository.UserRepository, secret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Username: req.Username,
		Password: string(hash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.TokenPair, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	tokenHash := hashToken(refreshToken)
	userID, err := s.repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	_ = s.repo.DeleteRefreshToken(ctx, tokenHash)
	return s.issueTokens(ctx, userID)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.repo.DeleteRefreshToken(ctx, tokenHash)
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	profile, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Profile{UserID: userID}, nil
		}
		return nil, err
	}
	return profile, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.Profile, error) {
	return s.repo.UpsertProfile(ctx, userID, req)
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*jwtutil.Claims, error) {
	return jwtutil.ValidateToken(tokenStr, s.jwtSecret)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uuid.UUID) (*model.TokenPair, error) {
	accessToken, err := jwtutil.GenerateAccessToken(userID, s.jwtSecret, s.accessTTL)
	if err != nil {
		return nil, err
	}
	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.repo.StoreRefreshToken(ctx, userID, hashToken(refreshToken), expiresAt); err != nil {
		return nil, err
	}
	return &model.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
