package authentication

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pr1/server/internal/app/authentication_service/auth_service"
	"pr1/server/internal/app/authentication_service/jwt_manager"
	"pr1/server/internal/app/authentication_service/users"
)

func Auth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	resp, err = handler(ctx, req)
	return resp, err
}

type AuthServer struct {
	auth_service.UnsafeAuthServiceServer
	userStore  users.UserStorage
	jwtManager *jwt_manager.JWTManager
}

func NewAuthServer(userStorage users.UserStorage, manager *jwt_manager.JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStorage,
		jwtManager: manager,
	}
}

func (a *AuthServer) Login(ctx context.Context, request *auth_service.LoginRequest) (*auth_service.LoginResponse, error) {
	user, err := a.userStore.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsCorrectPassword(request.Password) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	token, err := a.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token")
	}
	return &auth_service.LoginResponse{
		AccessToken: token,
	}, nil
}
