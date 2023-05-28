package client

import (
	"fmt"
	"log"
	"net/http"
	"pr1/client/internal/url_service_proto"
)

type Client interface {
	HandleRedirect(w http.ResponseWriter, r *http.Request)
	HandleCreation(w http.ResponseWriter, r *http.Request)
}

type client struct {
	grpc url_service_proto.ShortenerUrlClient
}

func NewClient(grpc url_service_proto.ShortenerUrlClient) Client {
	return &client{grpc: grpc}
}

func (c *client) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := r.URL.Path[1:]

	resp, err := c.grpc.GetShortenUrl(r.Context(), &url_service_proto.GetShortenUrlRequest{ShortenUrl: token})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, resp.Url, http.StatusMovedPermanently)
}

func (c *client) HandleCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := r.URL.Query().Get("url")

	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := c.grpc.CreateShortenUrl(r.Context(), &url_service_proto.CreateShortenUrlRequest{Url: url})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("http://%v/%v", r.Host, resp.ShortenUrl)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
