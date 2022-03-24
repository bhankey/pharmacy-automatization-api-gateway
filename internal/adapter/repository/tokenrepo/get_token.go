package tokenrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bhankey/go-utils/pkg/apperror"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (r *TokenRepo) GetToken(ctx context.Context, refreshToken string) (entities.RefreshToken, error) {
	errBase := fmt.Sprintf("tokenrepo.GetToken(%s)", refreshToken)

	const query = `
		SELECT id, user_id, refresh_token, user_agent, ip, finger_print, is_available, creation_time
		FROM refresh_tokens 
		WHERE refresh_token = $1 AND is_available = true
`

	var row struct {
		ID           int       `db:"id"`
		UserID       int       `db:"user_id"`
		RefreshToken string    `db:"refresh_token"`
		UserAgent    string    `db:"user_agent"`
		IP           string    `db:"ip"`
		FingerPrint  string    `db:"finger_print"`
		IsAvailable  bool      `db:"is_available"`
		CreationTime time.Time `db:"creation_time"`
	}

	if err := r.slave.GetContext(ctx, &row, query, refreshToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.RefreshToken{}, apperror.ErrNoEntity
		}

		return entities.RefreshToken{}, fmt.Errorf("%s: failed to get refresh token: %w", errBase, err)
	}

	return entities.RefreshToken{
		ID:           row.ID,
		UserID:       row.UserID,
		Token:        row.RefreshToken,
		UserAgent:    row.UserAgent,
		IP:           row.IP,
		FingerPrint:  row.FingerPrint,
		IsAvailable:  row.IsAvailable,
		CreationTime: row.CreationTime,
	}, nil
}
