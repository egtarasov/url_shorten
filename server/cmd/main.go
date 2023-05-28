package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"pr1/server/internal/app/authentication"
	"pr1/server/internal/app/repository"
	"pr1/server/internal/app/service"
	grpc_service "pr1/server/internal/app/url_service_proto"
)

func main() {
	repo := repository.NewUrlCacheRepo()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authentication.Auth),
	)
	implementation := service.NewService(repo)
	grpc_service.RegisterShortenerUrlServer(server, implementation)

	lsn, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("cant listen port 80")
	}

	if err = server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
