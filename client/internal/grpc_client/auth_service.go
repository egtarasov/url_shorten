package grpc_client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"pr1/client/internal/url_service_proto"
	"time"
)

type AuthService struct {
	service url_service_proto.AuthServiceClient
}

const (
	unauthorized         = "unauthorized"
	unauthorizedPassword = ""
)

func NewAuthService(ctx context.Context, target string) *AuthService {
	cc, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	service := url_service_proto.NewAuthServiceClient(cc)
	return &AuthService{
		service: service,
	}
}

func (a *AuthService) Login(username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &url_service_proto.LoginRequest{
		Username: username,
		Password: password,
	}
	resp, err := a.service.Login(ctx, req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return resp.AccessToken, err
}
