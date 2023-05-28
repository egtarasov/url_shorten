package repository

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type urlCacheRepo struct {
	cache *memcache.Client
}

func NewUrlCacheRepo() UrlRepo {
	cache := memcache.New("localhost:11211")
	return &urlCacheRepo{
		cache: cache,
	}
}

func (u urlCacheRepo) GetUrl(ctx context.Context, shortenUrl string) (*Url, error) {
	url, err := u.cache.Get(shortenUrl)
	if err != nil {
		return nil, err
	}
	return &Url{
		Url:        string(url.Value),
		ShortenUrl: url.Key,
		Token:      url.Key,
		// If the key is found in the cache, it means that the token has not expired yet.
		ExpirationTime: time.Now().Add(time.Minute),
	}, nil
}

func (u urlCacheRepo) CreateShortenUrl(ctx context.Context, url *Url) error {
	err := u.cache.Set(&memcache.Item{
		Key:        url.ShortenUrl,
		Value:      []byte(url.Url),
		Expiration: int32((time.Minute * 60).Seconds()),
	})
	return err
}

func (u urlCacheRepo) DeleteShortenUrl(ctx context.Context, shortenUrl string) error {
	err := u.cache.Delete(shortenUrl)
	return err
}
