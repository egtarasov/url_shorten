package jwt_manager

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"pr1/server/internal/app/authentication_service/users"
	"time"
)

type JWTManager interface {
	Verify(accessToken string) (*UserClaims, error)
	Create(user *users.User) (string, error)
}

type jwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJwtManager(secretKey string, tokenDuration time.Duration) JWTManager {
	return &jwtManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Verify verifies the access token and return user claims is the token is valid
func (m *jwtManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("invalid sign method")
			}
			return []byte(m.secretKey), nil
		})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %v", err)
	}
	return claims, nil
}

// Create generate access token based on user
func (m *jwtManager) Create(user *users.User) (string, error) {
	userClaims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return token.SignedString([]byte(m.secretKey))
}
