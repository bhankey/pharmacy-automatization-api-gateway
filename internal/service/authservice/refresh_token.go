package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhankey/go-utils/pkg/apperror"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
	identifyData entities.UserIdentifyData,
) (entities.Tokens, error) {
	errBase := fmt.Sprintf("authservice.RefreshToken(%s, %v)", refreshToken, identifyData)

	token, err := s.tokenStorage.GetToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, apperror.ErrNoEntity) {
			return entities.Tokens{}, apperror.NewClientError(apperror.WrongAuthToken, err)
		}

		return entities.Tokens{}, fmt.Errorf("%s: failed to get refresh token error: %w", errBase, err)
	}

	user, err := s.userStorage.GetByID(ctx, token.UserID)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("%s: failed to get user error: %w", errBase, err)
	}

	accessToken, err := s.createAndSignedToken(user.ID, user.Email, user.Role, user.DefaultPharmacyID, jwtExpireTime)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("failed to create access token error: %w", err)
	}

	newRefreshToken, err := s.createAndSaveRefreshToken(
		ctx,
		user.ID,
		user.Email,
		user.Role,
		user.DefaultPharmacyID,
		identifyData,
	)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("%s: failed to create refresh token error: %w", errBase, err)
	}

	go func() {
		ctx := context.Background()

		_ = s.tokenStorage.DeactivateTokenByIDs(ctx, []int{token.ID})
	}()

	return entities.Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
