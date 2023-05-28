package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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

type urlRepo struct {
	pool *pgxpool.Pool
}

func NewUrlRepo(pool *pgxpool.Pool) *urlRepo {
	return &urlRepo{pool: pool}
}

func (u *urlRepo) GetUrl(ctx context.Context, shortenUrl string) (*Url, error) {
	var url string
	var expirationTime time.Time
	err := u.pool.
		QueryRow(ctx, "SELECT Url, expire_at FROM public.urls WHERE shorten_url = $1", shortenUrl).
		Scan(&url, &expirationTime)
	if err != nil {
		return nil, err
	}
	return &Url{Url: url, ExpirationTime: expirationTime}, nil
}

func (u *urlRepo) CreateShortenUrl(ctx context.Context, url *Url) error {
	sql := `INSERT INTO public.urls 
    			(token, shorten_url, url, expired_at)
			VALUES 
			    ($1, $2, $3, $4)`
	_, err := u.pool.Exec(ctx, sql, url, url.Token, url.ShortenUrl, url.Url, url.ExpirationTime)
	return err
}

func (u *urlRepo) DeleteShortenUrl(ctx context.Context, shortenUrl string) error {
	_, err := u.pool.Exec(ctx, "DELETE FROM public.urls WHERE shorten_url = &1", shortenUrl)
	return err
}
