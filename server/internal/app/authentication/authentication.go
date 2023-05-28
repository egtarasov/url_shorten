package authentication

import (
	"context"
	"google.golang.org/grpc"
)

func Auth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	resp, err = handler(ctx, req)
	return resp, err
}
