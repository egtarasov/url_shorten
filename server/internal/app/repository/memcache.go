package repository

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type urlCacheRepo struct {
	cache *memcache.Client
}

func NewUrlCacheRepo() *urlCacheRepo {
	cache := memcache.New("localhost:11211")
	return &urlCacheRepo{
		cache: cache,
	}
}

func (u *urlCacheRepo) GetUrl(_ context.Context, token string) (*Url, error) {
	url, err := u.cache.Get(token)
	if err != nil {
		return nil, err
	}
	return &Url{
		Url:   string(url.Value),
		Token: url.Key,
		// If the key is found in the cache, it means that the token has not expired yet.
		ExpirationTime: time.Now().Add(time.Minute),
	}, nil
}

func (u *urlCacheRepo) CreateShortenUrl(_ context.Context, url *Url) error {
	err := u.cache.Set(&memcache.Item{
		Key:        url.Token,
		Value:      []byte(url.Url),
		Expiration: int32((time.Hour).Seconds()),
	})
	return err
}

func (u *urlCacheRepo) DeleteShortenUrl(_ context.Context, token string) error {
	err := u.cache.Delete(token)
	return err
}
