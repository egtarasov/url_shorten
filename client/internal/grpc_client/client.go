package grpc_client

import (
	"context"
	"google.golang.org/grpc"
	"pr1/client/internal/url_service_proto"
)

type Client struct {
	server url_service_proto.ShortenerUrlClient
	conn   *grpc.ClientConn
}

func NewClient(ctx context.Context, authInterceptor *AuthInterceptor, target string) (*Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))

	if err != nil {
		return nil, err
	}
	server := url_service_proto.NewShortenerUrlClient(conn)
	return &Client{server: server, conn: conn}, nil
}

func (c *Client) Close() {
	c.Close()
}

func (c *Client) CreateShortenUrl(ctx context.Context, in *url_service_proto.CreateShortenUrlRequest, opts ...grpc.CallOption) (*url_service_proto.CreateShortenUrlResponse, error) {
	resp, err := c.server.CreateShortenUrl(ctx, in)
	return resp, err
}

func (c *Client) GetShortenUrl(ctx context.Context, in *url_service_proto.GetShortenUrlRequest, opts ...grpc.CallOption) (*url_service_proto.GetShortenUrlResponse, error) {
	resp, err := c.server.GetShortenUrl(ctx, in)
	return resp, err
}
