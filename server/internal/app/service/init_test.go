package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"log"
	"pr1/server/internal/app/repository"
	"pr1/server/internal/app/repository/mocks"
	gen "pr1/server/internal/app/shorten_url_generator"
	"testing"
)

type ShortenerUrlServer struct {
	ctx      context.Context
	ctrl     *gomock.Controller
	repo     repository.UrlRepo
	repoMock *mocks.MockUrlRepo
	s        *service
}

func InitIntegration(t *testing.T) *ShortenerUrlServer {
	repo := repository.NewUrlCacheRepo()
	s := NewService(repo)
	return &ShortenerUrlServer{
		ctx:      context.Background(),
		repo:     repo,
		ctrl:     nil,
		repoMock: nil,
		s:        s,
	}
}

func (m *ShortenerUrlServer) AddTokens(urls []string) []string {
	tokens := make([]string, 0, len(urls))
	for _, url := range urls {
		token, err := gen.CreateShortenUrl(url)
		if err != nil {
			log.Fatal(err)
		}
		tokens = append(tokens, token)
		_ = m.repo.CreateShortenUrl(m.ctx, &repository.Url{Url: url, Token: token})
	}
	return tokens
}

func Init(t *testing.T) *ShortenerUrlServer {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockUrlRepo(ctrl)
	s := NewService(repo)
	return &ShortenerUrlServer{
		ctx:      context.Background(),
		ctrl:     ctrl,
		repoMock: repo,
		s:        s,
	}
}

func (m *ShortenerUrlServer) tearDown() {
	m.ctx.Done()
	if m.ctrl != nil {
		m.ctrl.Finish()
	}
}
