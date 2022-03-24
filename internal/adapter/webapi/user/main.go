package user

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	pb "github.com/bhankey/pharmacy-automatization-user/pkg/api/userservice"
)

type APIClient struct {
	client pb.UserServiceClient
}

func NewUserAPIClient(client pb.UserServiceClient) *APIClient {
	return &APIClient{
		client: client,
	}
}

func (c *APIClient) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	errBase := fmt.Sprintf("user.GetByEmail(%s)", email)

	user, err := c.client.GetByEmail(ctx, &pb.Email{
		Email: email,
	})
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", errBase, err)
	}

	return entities.User{
		ID:                int(user.Id),
		Name:              user.Name,
		Surname:           user.Surname,
		Email:             user.Email,
		Role:              entities.Role(user.Role),
		DefaultPharmacyID: int(user.DefaultPharmacyId),
	}, nil
}

func (c *APIClient) GetByID(ctx context.Context, id int) (entities.User, error) {
	errBase := fmt.Sprintf("user.GetByID(%d)", id)

	user, err := c.client.GetByID(ctx, &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", errBase, err)
	}

	return entities.User{
		ID:                int(user.Id),
		Name:              user.Name,
		Surname:           user.Surname,
		Email:             user.Email,
		Role:              entities.Role(user.Role),
		DefaultPharmacyID: int(user.DefaultPharmacyId),
	}, nil
}

func (c *APIClient) CreateUser(ctx context.Context, user entities.User) error {
	errBase := fmt.Sprintf("user.CreateUser(%v)", user)

	_, err := c.client.CreateUser(ctx, &pb.NewUser{
		Name:              user.Name,
		Email:             user.Email,
		Role:              string(user.Role),
		Password:          user.Password,
		UseIpCheck:        user.UseIPCheck,
		IsBlocked:         user.IsBlock,
		DefaultPharmacyId: int64(user.DefaultPharmacyID),
		Surname:           user.Surname,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", errBase, err)
	}

	return nil
}

func (c *APIClient) RequestToChangePassword(ctx context.Context, email string) error {
	errBase := fmt.Sprintf("user.RequestToChangePassword(%s)", email)

	_, err := c.client.RequestToChangePassword(ctx, &pb.Email{
		Email: email,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", errBase, err)
	}

	return nil
}

func (c *APIClient) ChangePassword(ctx context.Context, email, code, newPassword string) error {
	errBase := fmt.Sprintf("user.ChangePassword(%s, %s, %s)", email, code, newPassword)

	_, err := c.client.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Email:       email,
		Code:        code,
		NewPassword: newPassword,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", errBase, err)
	}

	return nil
}

func (c *APIClient) GetUsers(ctx context.Context, lastID, limit int) ([]entities.User, error) {
	errBase := fmt.Sprintf("user.GetUsers(%d, %d)", lastID, limit)

	users, err := c.client.GetUsers(ctx, &pb.PaginationRequest{
		LastId: int64(lastID),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errBase, err)
	}

	res := make([]entities.User, 0, len(users.GetUsers()))
	for _, user := range users.GetUsers() {
		res = append(res, entities.User{
			ID:                int(user.Id),
			Name:              user.Name,
			Surname:           user.Surname,
			Email:             user.Email,
			Role:              entities.Role(user.Role),
			UseIPCheck:        user.UseIpCheck,
			IsBlock:           user.IsBlocked,
			DefaultPharmacyID: int(user.DefaultPharmacyId),
		})
	}

	return res, nil
}

func (c *APIClient) UpdateUser(ctx context.Context, user entities.User) error {
	errBase := fmt.Sprintf("user.UpdateUser(%v)", user)

	_, err := c.client.UpdateUser(ctx, &pb.User{
		Id:                int64(user.ID),
		Name:              user.Name,
		Email:             user.Email,
		Role:              string(user.Role),
		UseIpCheck:        user.UseIPCheck,
		IsBlocked:         user.IsBlock,
		DefaultPharmacyId: int64(user.DefaultPharmacyID),
		Surname:           user.Surname,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", errBase, err)
	}

	return nil
}

func (c *APIClient) IsPasswordCorrect(ctx context.Context, email, password string) (bool, error) {
	errBase := fmt.Sprintf("user.IsPasswordCorrect(%s, %s)", email, password)

	isCorrect, err := c.client.IsPasswordCorrect(ctx, &pb.EmailAndPassword{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", errBase, err)
	}

	return isCorrect.IsCorrect, nil
}
