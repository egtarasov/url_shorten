package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
}

func NewUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	return &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
	}, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
