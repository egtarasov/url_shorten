package repository

import (
	"context"
	"time"
)

type Url struct {
	Url            string
	ShortenUrl     string
	Token          string
	ExpirationTime time.Time
}

type UrlRepo interface {
	GetUrl(ctx context.Context, shortenUrl string) (*Url, error)
	CreateShortenUrl(ctx context.Context, url *Url) error
	DeleteShortenUrl(ctx context.Context, shortenUrl string) error
}
