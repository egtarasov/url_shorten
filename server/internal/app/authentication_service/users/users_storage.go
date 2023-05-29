package users

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserStorage interface {
	// GetByUsername returns a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)
	// Add adds a new user
	Add(ctx context.Context, user *User) error
	// Update updates a user
	Update(ctx context.Context, user *User) error
}

type userStorage struct {
	repo *pgxpool.Pool
}

func NewUserStorage(repo *pgxpool.Pool) UserStorage {
	return &userStorage{
		repo: repo,
	}
}

func (u *userStorage) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user *User
	err := pgxscan.Get(ctx, u.repo, user, "SELECT username, hashed_password FROM users WHERE username = $1", username)
	return user, err
}

func (u *userStorage) Add(ctx context.Context, user *User) error {
	_, err := u.repo.Exec(ctx, "INSERT INTO users (username, hashed_password) VALUES ($1, $2)", user.Username, user.HashedPassword)
	return err
}

func (u *userStorage) Update(ctx context.Context, user *User) error {
	_, err := u.repo.Exec(ctx, "UPDATE users SET hashed_password = $1 WHERE username = $2", user.HashedPassword, user.Username)
	return err
}
