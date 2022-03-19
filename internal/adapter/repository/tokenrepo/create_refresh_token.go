package tokenrepo

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (r *TokenRepo) CreateRefreshToken(ctx context.Context, token entities.RefreshToken) error {
	errBase := fmt.Sprintf("tokenrepo.CreateRefreshToken(%v)", token)

	const query = `
		INSERT INTO refresh_tokens(user_id, refresh_token, user_agent, ip, finger_print)
							VALUES ($1, $2, $3, $4, $5)
`

	if _, err := r.master.ExecContext(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.UserAgent,
		token.IP,
		token.FingerPrint,
	); err != nil {
		return fmt.Errorf("%s: failed to insert refresh token: %w", errBase, err)
	}

	return nil
}
