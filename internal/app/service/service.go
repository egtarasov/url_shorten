package service

import (
	"context"
	"errors"
	"pr1/internal/app/repository"
	gen "pr1/internal/app/shorten_url_generator"
	"pr1/internal/app/url_service_proto"
	"time"
)

type service struct {
	url_service_proto.UnimplementedShortenerUrlServer
	repo repository.UrlRepo
}

func (s *service) CreateShortenUrl(ctx context.Context, request *url_service_proto.CreateShortenUrlRequest) (*url_service_proto.CreateShortenUrlResponse, error) {
	shortUrl, err := gen.CreateShortenUrl()
	if err != nil {
		return nil, err
	}
	url := &repository.Url{
		Url:            request.Url,
		ShortenUrl:     shortUrl.UrlShorten,
		Token:          shortUrl.Token,
		ExpirationTime: time.Now().Add(time.Hour * 24),
	}
	err = s.repo.CreateShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	return &url_service_proto.CreateShortenUrlResponse{ShortenUrl: url.ShortenUrl}, nil
}

func (s *service) GetShortenUrl(ctx context.Context, request *url_service_proto.GetShortenUrlRequest) (*url_service_proto.GetShortenUrlResponse, error) {
	url, err := s.repo.GetUrl(ctx, request.ShortenUrl)
	if err != nil {
		return nil, err
	}
	if url.ExpirationTime.After(time.Now()) {
		return nil, errors.New("url is expired")
	}

	return &url_service_proto.GetShortenUrlResponse{Url: url.Url}, nil
}
