package auth_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"pr1/server/internal/app/authentication_service/jwt_manager"
)

type AuthInterceptor struct {
	accessibleRoles map[string][]string
	jwtManager      jwt_manager.JWTManager
}

func Roles() map[string][]string {
	return map[string][]string{
		"/ShortenerUrl/CreateShortenUrl": {"user", "admin"},
	}
}

func NewAuthInterceptor(jwtManager jwt_manager.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{
		accessibleRoles: Roles(),
		jwtManager:      jwtManager,
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

		err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) error {
	roles, ok := i.accessibleRoles[method]
	if !ok {
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.NotFound, "metadata is not provided")
	}
	values, ok := md["authorization"]
	if len(values) == 0 {
		return status.Error(codes.NotFound, "authorization token is not provided")
	}
	accessToken := values[0]
	userClaims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, v := range roles {
		if v == userClaims.Role {
			return nil
		}
	}
	return status.Error(codes.PermissionDenied, "no permission granted")
}
