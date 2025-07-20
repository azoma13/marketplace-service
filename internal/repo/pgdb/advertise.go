package pgdb

import (
	"context"
	"fmt"

	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/pkg/postgres"
)

type AdvertiseRepo struct {
	*postgres.Postgres
}

func NewAdvertiseRepo(pg *postgres.Postgres) *AdvertiseRepo {
	return &AdvertiseRepo{pg}
}

func (r *AdvertiseRepo) CreateAdvertise(ctx context.Context, advertise entity.Advertise) (entity.Advertise, error) {
	query := `
	INSERT INTO advertises (title, description, image, price, user_id)
		VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at;`

	err := r.Pool.QueryRow(ctx, query, advertise.Title, advertise.Description, advertise.Image, advertise.Price, advertise.UserId).Scan(&advertise.Id, &advertise.CreatedAt)
	if err != nil {
		return entity.Advertise{}, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %w", err)
	}

	return advertise, nil
}

type AdvertiseGetFeedAdRepo struct {
	UserId    int
	Sort      string
	SortOrder string
	Page      int
	PerPage   int
	MinPrice  float64
	MaxPrice  float64
}

func (r *AdvertiseRepo) GetFeedAdvertise(ctx context.Context, input AdvertiseGetFeedAdRepo) ([]*entity.FeedAdvertises, error) {
	var offset int
	if input.Page > 0 {
		offset = (input.Page - 1) * input.PerPage
	} else {
		input.Page = 1
		offset = 0
	}

	query := `
	SELECT
    	a.title,
    	a.description,
    	a.image,
    	a.price,
    	u.username AS author_username,
    	CASE WHEN a.user_id = $1 THEN TRUE ELSE FALSE END AS is_author
	FROM
    	advertises a
	JOIN
    	users u ON a.user_id = u.id
	WHERE
    	(a.price BETWEEN $2 AND $3)
	ORDER BY
    	` + input.Sort + ` ` + input.SortOrder + `
	LIMIT
    	$4 OFFSET $5;
	`

	rows, err := r.Pool.Query(ctx, query, input.UserId, input.MinPrice, input.MaxPrice, input.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("error query for get ads: %v", err)
	}
	defer rows.Close()

	advertises := []*entity.FeedAdvertises{}
	for rows.Next() {
		var advertise entity.FeedAdvertises
		err := rows.Scan(&advertise.Title, &advertise.Description, &advertise.ImageUrl, &advertise.Price, &advertise.AuthorUsername, &advertise.IsAuthor)
		if err != nil {
			return nil, fmt.Errorf("error scan rows for get ad: %v", err)
		}
		advertises = append(advertises, &advertise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error cursor rows in func GetFeedAdvertise: %v", err)
	}

	return advertises, nil
}
