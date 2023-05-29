package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pr1/server/internal/app/shorten_url_generator"
	"pr1/server/internal/app/url_service_proto"
	"testing"
)

func TestService_CreateShortenUrl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		m := Init(t)
		defer m.tearDown()

		urls := []string{"https://hostiq.ua/blog/what-is-url/",
			"https://developer.mozilla.org/ru/docs/Learn/Common_questions/Web_mechanics/What_is_a_URL",
			"https://www.qrcode-tiger.com/ru/create-multiple-unique-qr-code-for-url",
			"https://support.ariba.com/item/view/KB0393053_ru?HelpCenter=1}",
		}

		for _, url := range urls[:1] {
			m.repoMock.EXPECT().CreateShortenUrl(m.ctx, gomock.Any()).Return(nil)
			resp, err := m.s.CreateShortenUrl(m.ctx, &url_service_proto.CreateShortenUrlRequest{Url: url})
			require.NoError(t, err)

			expect, err := shorten_url_generator.CreateShortenUrl(url)
			require.NoError(t, err)
			assert.Equal(t, expect, resp.ShortenUrl)
		}
	})
}

func TestService_GetShortenUrl(t *testing.T) {
	t.Parallel()
	m := InitIntegration(t)
	defer m.tearDown()

	tt := struct {
		urls             []string
		urlsWithoutToken []string
		tokens           []string
	}{
		urls:             []string{"url1", "url2", "url3"},
		urlsWithoutToken: []string{"url2", "ursdfl2", "url4"},
		tokens:           m.AddTokens([]string{"url1", "url2", "url3"}),
	}

	for i, expectUrl := range tt.urls {
		url, err := m.s.GetShortenUrl(m.ctx, &url_service_proto.GetShortenUrlRequest{ShortenUrl: tt.tokens[i]})
		require.NoError(t, err)
		assert.Equal(t, expectUrl, url.Url)
	}
	for _, url := range tt.urlsWithoutToken {
		_, err := m.s.GetShortenUrl(m.ctx, &url_service_proto.GetShortenUrlRequest{ShortenUrl: url})
		require.Error(t, err)
	}
}
