package auth_service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"pr1/server/internal/app/authentication_service/jwt_manager"
	"pr1/server/internal/app/authentication_service/users"
	url_service_proto2 "pr1/server/internal/app/url_service_proto"
)

type AuthService struct {
	url_service_proto2.UnimplementedAuthServiceServer
	storage    users.UserStorage
	jwtManager jwt_manager.JWTManager
}

func NewAuthService(storage users.UserStorage, jwtManager jwt_manager.JWTManager) *AuthService {
	return &AuthService{
		storage:    storage,
		jwtManager: jwtManager,
	}
}

func (a *AuthService) Login(ctx context.Context, request *url_service_proto2.LoginRequest) (*url_service_proto2.LoginResponse, error) {
	log.Println("--> Login[AuthService]: ", request.Username)
	user, err := a.storage.GetByUsername(ctx, request.Username)
	if err != nil {
		log.Printf("error:%v\n", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	if user == nil || !user.IsCorrectPassword(request.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect username/password")
	}
	token, err := a.jwtManager.Create(user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cant create access token")
	}
	return &url_service_proto2.LoginResponse{AccessToken: token}, nil
}
