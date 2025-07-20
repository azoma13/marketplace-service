package service

import (
	"context"
	"time"

	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/internal/repo"
	"github.com/azoma13/marketplace-service/pkg/hasher"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	ParseToken(accessToken string) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
}

type AdvertiseCreateNewAdvertiseInput struct {
	Title       string
	Description string
	Image       string
	Price       float64
	UserId      int
}

type AdvertiseDowlandImageInput struct {
	ImageUrl string
	UserId   int
}

type AdvertiseGetFeedAdInput struct {
	Sort         string
	Page         string
	PerPage      string
	CurrentPrice string
	UserId       int
}

type Advertise interface {
	CreateAdvertise(ctx context.Context, input AdvertiseCreateNewAdvertiseInput) (entity.Advertise, error)
	GetFeedAdvertise(ctx context.Context, input AdvertiseGetFeedAdInput) ([]*entity.FeedAdvertises, error)
}

type FeedAdvertises interface{}

type Services struct {
	Auth
	Advertise
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewService(deps ServicesDependencies) *Services {
	return &Services{
		Auth:      NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Advertise: NewAdvertiseService(deps.Repos.Advertise, deps.Repos.User),
	}
}
