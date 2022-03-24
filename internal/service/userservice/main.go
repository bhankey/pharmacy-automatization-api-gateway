package userservice

import (
	"context"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

type UserService struct {
	userStorage UserStorage
}

type UserStorage interface {
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	GetByID(ctx context.Context, id int) (entities.User, error)
	CreateUser(ctx context.Context, user entities.User) error
	RequestToChangePassword(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, email, code, newPassword string) error
	GetUsers(ctx context.Context, lastID, limit int) ([]entities.User, error)
	UpdateUser(ctx context.Context, user entities.User) error
}

func NewUserService(
	userStorage UserStorage,
) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}
