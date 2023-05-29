package main

import (
	"context"
	"log"
	"net/http"
	"pr1/client/internal/client"
	"pr1/client/internal/grpc_client"
)

const (
	target = "localhost:80"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	authService := grpc_client.NewAuthService(ctx, target)

	cli := client.NewClient(ctx, target, authService)
	mux := http.NewServeMux()

	mux.HandleFunc("/create", cli.HandleCreation)
	mux.HandleFunc("/", cli.HandleRedirect)

	log.Println("starting server at :1200")
	if err := http.ListenAndServe(":1200", mux); err != nil {
		log.Fatal(err)
	}
}
