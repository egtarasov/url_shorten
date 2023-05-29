//go:generate mockgen -source=url_repo.go -destination=mocks/url_repo_mock.go -package=mocks
package repository

import (
	"context"
	"time"
)

type Url struct {
	Url            string
	Token          string
	ExpirationTime time.Time
}

type UrlRepo interface {
	GetUrl(ctx context.Context, token string) (*Url, error)
	CreateShortenUrl(ctx context.Context, url *Url) error
	DeleteShortenUrl(ctx context.Context, token string) error
}
