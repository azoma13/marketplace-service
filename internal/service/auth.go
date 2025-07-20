package service

import (
	"context"
	"fmt"
	"time"

	"github.com/azoma13/marketplace-service/config"
	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/internal/repo"
	"github.com/azoma13/marketplace-service/pkg/hasher"
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error) {
	user := &entity.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.userRepo.CreateUser(ctx, *user)
	if err != nil {
		return 0, fmt.Errorf("AuthService.CreateUser - s.userRepo.CreateUser: %w", err)
	}

	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, s.passwordHasher.Hash(input.Password))
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateToken - s.userRepo.GetUserByUsernameAndPassword: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "token",
		},
		UserId: user.Id,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateToken - token.SignedString: %v", err)
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Cfg.SignKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parse token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("error parse token")
	}

	return claims.UserId, nil
}
