package repo

import (
	"context"

	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/internal/repo/pgdb"

	"github.com/azoma13/marketplace-service/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type Advertise interface {
	CreateAdvertise(ctx context.Context, advertise entity.Advertise) (entity.Advertise, error)
	GetFeedAdvertise(ctx context.Context, input pgdb.AdvertiseGetFeedAdRepo) ([]*entity.FeedAdvertises, error)
}

type Repositories struct {
	User
	Advertise
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:      pgdb.NewUserRepo(pg),
		Advertise: pgdb.NewAdvertiseRepo(pg),
	}
}
