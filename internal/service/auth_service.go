package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/BernadDwiki/shortlink-backend/internal/dto"
	"github.com/BernadDwiki/shortlink-backend/internal/jwt"
	"github.com/BernadDwiki/shortlink-backend/internal/model"
	"github.com/BernadDwiki/shortlink-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type AuthService struct {
	repo   *repository.UserRepository
	secret string
	ttl    time.Duration
}

func NewAuthService(repo *repository.UserRepository, secret string, ttl time.Duration) *AuthService {
	return &AuthService{repo: repo, secret: secret, ttl: ttl}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyRegistered
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, req.Email, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, s.secret, s.ttl)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
