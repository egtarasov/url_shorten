package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"pr1/server/internal/app/authentication_service/auth_interceptor"
	"pr1/server/internal/app/authentication_service/users"
	"pr1/server/internal/app/repository"
	"pr1/server/internal/app/service"
	grpc_service "pr1/server/internal/app/url_service_proto"
)

const addUsers = false

func AddUsers(ctx context.Context, storage users.UserStorage) {
	usersAdd := []struct {
		username string
		password string
	}{
		{"bubon", "12345"},
		{"bob", "12345"},
		{"buldozer", "12345"},
		{"abc", "12345"},
	}
	for _, user := range usersAdd {
		_ = storage.Add(ctx, user.username, user.password)
	}
}

func connectionString() string {
	return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("can't load environment variables")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.Connect(ctx, connectionString())
	if err != nil {
		log.Fatal(fmt.Sprintf("can't connect ot db: %v", err))
	}
	userStorage := users.NewUserStorage(pool)

	authInterceptor := auth_interceptor.NewAuthInterceptor(userStorage)

	repo := repository.NewUrlCacheRepo()

	urlShorterService := service.NewService(repo)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	grpc_service.RegisterShortenerUrlServer(server, urlShorterService)

	lsn, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("cant listen port 80")
	}

	// Add some users
	if addUsers {
		AddUsers(ctx, userStorage)
	}

	if err = server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
