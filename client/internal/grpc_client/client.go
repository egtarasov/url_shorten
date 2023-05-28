package grpc_client

import (
	"context"
	"google.golang.org/grpc"
	"pr1/client/internal/url_service_proto"
)

type client struct {
	server url_service_proto.ShortenerUrlClient
	conn   *grpc.ClientConn
}

func NewClient(ctx context.Context, target string) (*client, error) {
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	server := url_service_proto.NewShortenerUrlClient(conn)
	return &client{server: server, conn: conn}, nil
}

func (c *client) Close() {
	c.Close()
}

func (c *client) CreateShortenUrl(ctx context.Context, in *url_service_proto.CreateShortenUrlRequest, opts ...grpc.CallOption) (*url_service_proto.CreateShortenUrlResponse, error) {
	resp, err := c.server.CreateShortenUrl(ctx, in)
	return resp, err
}

func (c *client) GetShortenUrl(ctx context.Context, in *url_service_proto.GetShortenUrlRequest, opts ...grpc.CallOption) (*url_service_proto.GetShortenUrlResponse, error) {
	resp, err := c.server.GetShortenUrl(ctx, in)
	return resp, err
}
