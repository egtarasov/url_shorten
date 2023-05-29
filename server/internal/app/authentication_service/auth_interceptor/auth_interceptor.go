package auth_interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"pr1/server/internal/app/authentication_service/users"
	url_service_proto2 "pr1/server/internal/app/url_service_proto"
)

type AuthInterceptor struct {
	storage users.UserStorage
}

func NewAuthInterceptor(storage users.UserStorage) *AuthInterceptor {
	return &AuthInterceptor{
		storage: storage,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		// Lack of 'roles' for user in authorization is compensated by hard coding
		if info.FullMethod == "/ShortenerUrl/GetShortenUrl" {
			return handler(ctx, req)
		}

		request, ok := req.(*url_service_proto2.CreateShortenUrlRequest)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid request")
		}
		user, err := i.storage.GetByUsername(ctx, request.UserName)
		if err != nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("cant find user in db: %v", err))
		}
		if !user.IsCorrectPassword(request.Password) {
			return nil, status.Error(codes.Unauthenticated, "invalid username/password")
		}
		return handler(ctx, req)
	}
}
