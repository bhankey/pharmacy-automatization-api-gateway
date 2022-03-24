package userservice

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (s *UserService) UpdateUser(ctx context.Context, user entities.User) error {
	errBase := fmt.Sprintf("userservice.UpdateUser(%v)", user)

	if err := s.userStorage.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("%s: failed to update user: %w", errBase, err)
	}

	return nil
}

func (s *UserService) GetBatchOfUsers(ctx context.Context, lastClientID int, limit int) ([]entities.User, error) {
	errBase := fmt.Sprintf("userservice.GetBatchOfUsers(%d, %d)", lastClientID, limit)

	users, err := s.userStorage.GetUsers(ctx, lastClientID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get users: %w", errBase, err)
	}

	return users, nil
}

func (s *UserService) Registry(ctx context.Context, user entities.User) error {
	errBase := fmt.Sprintf("userservice.Registry(%v)", user)

	if err := s.userStorage.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("%s: failed to create user error: %w", errBase, err)
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	errBase := fmt.Sprintf("userservice.ResetPassword(%s, %s, %s)", email, code, newPassword)

	if err := s.userStorage.ChangePassword(ctx, email, code, newPassword); err != nil {
		return fmt.Errorf("%s failed to change password: %w", errBase, err)
	}

	return nil
}

func (s *UserService) RequestToResetPassword(ctx context.Context, email string) error {
	errBase := fmt.Sprintf("userservice.RequestToResetPassword(%s)", email)

	if err := s.userStorage.RequestToChangePassword(ctx, email); err != nil {
		return fmt.Errorf("%s: failed to request password change: %w", errBase, err)
	}

	return nil
}
