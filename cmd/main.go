package main

import (
	"google.golang.org/grpc"
	"pr1/internal/app/authentication"
)

func main() {
	_ = grpc.NewServer(
		grpc.UnaryInterceptor(authentication.Auth),
	)

}
