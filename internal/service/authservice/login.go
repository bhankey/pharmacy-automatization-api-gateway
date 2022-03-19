package authservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/apperror"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"sort"
	"time"
)

const (
	jwtExpireTime        = time.Minute * 15
	jwtExpireRefreshTime = time.Hour * 24 * 60
	maxActivateToken     = 5
	slaveLag             = 5
)

// Move this to some service that implement registry pattern (consul).
func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	identifyData entities.UserIdentifyData,
) (entities.Tokens, error) {
	errBase := fmt.Sprintf("authservice.Login(%s, %s, %v)", email, password, identifyData)

	isCorrect, err := s.userStorage.IsPasswordCorrect(ctx, email, password)
	if err != nil && !errors.Is(err, apperror.ErrNoEntity) {
		return entities.Tokens{}, fmt.Errorf("%s :failed to check password: %w", errBase, err)
	}

	if !isCorrect {
		return entities.Tokens{}, apperror.NewClientError(apperror.WrongAuthorization, err)
	}

	user, err := s.userStorage.GetByEmail(ctx, email)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("%s :failed to get user: %w", errBase, err)
	}

	accessToken, err := s.createAndSignedToken(user.ID, user.Email, user.Role, user.DefaultPharmacyID, jwtExpireTime)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("%s: failed to create and sigend accass token error: %w", errBase, err)
	}

	refreshToken, err := s.createAndSaveRefreshToken(
		ctx,
		user.ID,
		user.Email,
		user.Role,
		user.DefaultPharmacyID,
		identifyData,
	)
	if err != nil {
		return entities.Tokens{}, fmt.Errorf("%s: failed to create and sigend refresh token error: %w", errBase, err)
	}

	// TODO write wrap to catch panic in goroutine
	go func() {
		ctx := context.Background()

		time.Sleep(time.Second * slaveLag) // In 20 seconds slave will surely update data

		_ = s.deactivateMaxReachedTokensCount(ctx, user.ID) // TODO think about async error handling
	}()

	return entities.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) deactivateMaxReachedTokensCount(ctx context.Context, userID int) error {
	tokens, err := s.tokenStorage.GetAllActiveRefreshTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get all active tokens error: %w", err)
	}

	if len(tokens) > maxActivateToken {
		tokensIDsToDeativate := make([]int, 0, len(tokens)-maxActivateToken)
		sort.Slice(tokens, func(i, j int) bool {
			return tokens[i].ID < tokens[j].ID
		})

		for i := 0; i < len(tokens)-maxActivateToken; i++ {
			tokensIDsToDeativate = append(tokensIDsToDeativate, tokens[i].ID)
		}

		if err := s.tokenStorage.DeactivateTokenByIDs(ctx, tokensIDsToDeativate); err != nil {
			return fmt.Errorf("failed to deativate old tokens error: %w", err)
		}
	}

	return nil
}
