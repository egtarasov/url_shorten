package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"pr1/client/internal/grpc_client"
	"pr1/client/internal/url_service_proto"
	"time"
)

var (
	refreshDuration = time.Minute * 5
)

type Client interface {
	HandleRedirect(w http.ResponseWriter, r *http.Request)
	HandleCreation(w http.ResponseWriter, r *http.Request)
}

type client struct {
	grpc        url_service_proto.ShortenerUrlClient
	interceptor *grpc_client.AuthInterceptor
}

func NewClient(ctx context.Context, target string, authService *grpc_client.AuthService) Client {
	interceptor := grpc_client.NewAuthInterceptor(authService, refreshDuration)
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()))

	if err != nil {
		log.Fatal(err)
	}
	server := url_service_proto.NewShortenerUrlClient(conn)
	return &client{grpc: server, interceptor: interceptor}
}

func (c *client) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleRedirect ->")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid method")
		return
	}
	token := r.URL.Path[1:]

	resp, err := c.grpc.GetShortenUrl(r.Context(), &url_service_proto.GetShortenUrlRequest{ShortenUrl: token})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	http.Redirect(w, r, resp.Url, http.StatusMovedPermanently)
}

func (c *client) HandleCreation(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleCreation ->")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid method")
		return
	}

	url := r.URL.Query().Get("url")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if url == "" || username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("no username, url or password")
		return
	}
	c.interceptor.SetCredentials(username, password)

	resp, err := c.grpc.CreateShortenUrl(r.Context(), &url_service_proto.CreateShortenUrlRequest{
		Url:      url,
		Username: username,
		Password: password,
	})
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
