package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

type urlPostgresRepo struct {
	pool *pgxpool.Pool
}

func connectionString() string {
	return fmt.Sprintf("host=localhost port=5432 user=%s password=%s db=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
}

func NewUrlPostgresRepo(ctx context.Context) UrlRepo {
	pool, err := pgxpool.Connect(ctx, connectionString())
	if err != nil {
		log.Fatal("cant connect to db")
	}
	return &urlPostgresRepo{pool: pool}
}

func (u *urlPostgresRepo) GetUrl(ctx context.Context, shortenUrl string) (*Url, error) {
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

func (u *urlPostgresRepo) CreateShortenUrl(ctx context.Context, url *Url) error {
	sql := `INSERT INTO public.urls 
    			(token, shorten_url, url, expired_at)
			VALUES 
			    ($1, $2, $3, $4)`
	_, err := u.pool.Exec(ctx, sql, url, url.Token, url.ShortenUrl, url.Url, url.ExpirationTime)
	return err
}

func (u *urlPostgresRepo) DeleteShortenUrl(ctx context.Context, shortenUrl string) error {
	_, err := u.pool.Exec(ctx, "DELETE FROM public.urls WHERE shorten_url = &1", shortenUrl)
	return err
}
