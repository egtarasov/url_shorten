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
	Add(ctx context.Context, username, password, role string) error
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
	var user User
	err := pgxscan.Get(ctx, u.repo, &user, "SELECT username, hashed_password, role FROM users WHERE username = $1", username)
	return &user, err
}

func (u *userStorage) Add(ctx context.Context, username, password, role string) error {
	user, err := NewUser(username, password, role)
	if err != nil {
		return err
	}
	_, err = u.repo.Exec(ctx, "INSERT INTO users (username, hashed_password, role) VALUES ($1, $2, $3)",
		user.Username, user.HashedPassword, user.Role)
	return err
}

func (u *userStorage) Update(ctx context.Context, user *User) error {
	_, err := u.repo.Exec(ctx, "UPDATE users SET hashed_password = $1, role = $2 WHERE username = $3", user.HashedPassword, user.Role, user.Username)
	return err
}
