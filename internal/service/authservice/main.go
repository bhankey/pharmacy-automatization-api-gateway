package authservice

import (
	"context"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

type AuthService struct {
	userStorage  UserStorage
	tokenStorage TokenStorage

	jwtKey string
}

type UserStorage interface {
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	GetByID(ctx context.Context, id int) (entities.User, error)
	IsPasswordCorrect(ctx context.Context, email, password string) (bool, error)
}

type TokenStorage interface {
	CreateRefreshToken(ctx context.Context, token entities.RefreshToken) error
	GetAllActiveRefreshTokens(ctx context.Context, userID int) ([]entities.RefreshToken, error)
	DeactivateTokenByIDs(ctx context.Context, tokenIDs []int) error
	GetToken(ctx context.Context, refreshToken string) (entities.RefreshToken, error)
}

func NewAuthService(
	userStorage UserStorage,
	tokenStorage TokenStorage,
	jwtKey string,
) *AuthService {
	return &AuthService{
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		jwtKey:       jwtKey,
	}
}
