package grpc_client

import (
	"context"
	"golang.org/x/crypto/openpgp/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type AuthInterceptor struct {
	authService *AuthService
	authMethods map[string]bool
	username    string
	password    string
	accessToken string
}

func authMethods() map[string]bool {
	return map[string]bool{
		"/ShortenerUrl/GetShortenUrl":    false,
		"/ShortenerUrl/CreateShortenUrl": true,
	}
}

func NewAuthInterceptor(authService *AuthService, refreshDuration time.Duration) *AuthInterceptor {
	interceptor := &AuthInterceptor{
		authService: authService,
		authMethods: authMethods(),
		username:    unauthorized,
		password:    unauthorizedPassword,
	}
	err := interceptor.scheduleTokenRefreshment(refreshDuration)
	if err != nil {
		log.Fatal(err)
	}
	return interceptor
}

func (i *AuthInterceptor) Unary() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("intercepted method: ", method)

		if !i.authMethods[method] {
			// skip authentication for this method
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.accessToken)
}

func (i *AuthInterceptor) refreshToken() error {
	accessToken, err := i.authService.Login(i.username, i.password)
	if err != nil {
		return err
	}
	i.accessToken = accessToken
	log.Printf("token refreshed: %v\n", accessToken)
	return nil
}

func (i *AuthInterceptor) scheduleTokenRefreshment(refreshDuration time.Duration) error {
	err := i.refreshToken()
	if err != nil {
		return err
	}
	go func(wait time.Duration) {
		for i.username == unauthorized {
			time.Sleep(time.Second * 5)
			log.Println(i.accessToken)
			log.Println("waiting for credentials to be set")
		}
		for {
			err := i.refreshToken()
			if err != nil {
				// wait for a second to retry if the refreshment failed
				wait = time.Second
			} else {
				wait = refreshDuration
			}
			time.Sleep(wait)
		}
	}(refreshDuration)

	return nil
}

func (i *AuthInterceptor) SetCredentials(username, password string) {
	i.username = username
	i.password = password
	err := errors.ErrUnknownIssuer
	for err != nil {
		err = i.refreshToken()
	}
}
