package service

import (
	"context"
	"errors"
	"pr1/server/internal/app/repository"
	gen "pr1/server/internal/app/shorten_url_generator"
	url_service_proto2 "pr1/server/internal/app/url_service_proto"
	"time"
)

type service struct {
	url_service_proto2.UnimplementedShortenerUrlServer
	repo repository.UrlRepo
}

func NewService(repo repository.UrlRepo) *service {
	return &service{repo: repo}
}

func (s *service) CreateShortenUrl(ctx context.Context, request *url_service_proto2.CreateShortenUrlRequest) (*url_service_proto2.CreateShortenUrlResponse, error) {
	token, err := gen.CreateShortenUrl(request.Url)
	if err != nil {
		return nil, err
	}
	url := &repository.Url{
		Url:            request.Url,
		Token:          token,
		ExpirationTime: time.Now().Add(time.Hour * 24),
	}
	err = s.repo.CreateShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	return &url_service_proto2.CreateShortenUrlResponse{ShortenUrl: token}, nil
}

func (s *service) GetShortenUrl(ctx context.Context, request *url_service_proto2.GetShortenUrlRequest) (*url_service_proto2.GetShortenUrlResponse, error) {
	url, err := s.repo.GetUrl(ctx, request.ShortenUrl)
	if err != nil {
		return nil, err
	}
	if !url.ExpirationTime.After(time.Now()) {
		return nil, errors.New("url is expired")
	}

	return &url_service_proto2.GetShortenUrlResponse{Url: url.Url}, nil
}
