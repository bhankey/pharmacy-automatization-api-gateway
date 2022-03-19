package tokenrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (r *TokenRepo) GetAllActiveRefreshTokens(ctx context.Context, userID int) ([]entities.RefreshToken, error) {
	errBase := fmt.Sprintf("tokenrepo.GetAllActiveRefreshTokens(%d)", userID)

	const query = `
		SELECT id, user_id, refresh_token, user_agent, ip, finger_print, is_available, creation_time
		FROM refresh_tokens 
		WHERE user_id = $1 AND is_available = true
`

	var rows []struct {
		ID           int       `db:"id"`
		UserID       int       `db:"user_id"`
		RefreshToken string    `db:"refresh_token"`
		UserAgent    string    `db:"user_agent"`
		IP           string    `db:"ip"`
		FingerPrint  string    `db:"finger_print"`
		IsAvailable  bool      `db:"is_available"`
		CreationTime time.Time `db:"creation_time"`
	}

	if err := r.slave.SelectContext(ctx, &rows, query, userID); err != nil {
		return nil, fmt.Errorf("%s: failed to get all refresh tokens: %w", errBase, err)
	}

	result := make([]entities.RefreshToken, 0, len(rows))
	for _, row := range rows {
		result = append(result, entities.RefreshToken{
			ID:           row.ID,
			UserID:       row.UserID,
			Token:        row.RefreshToken,
			UserAgent:    row.UserAgent,
			IP:           row.IP,
			FingerPrint:  row.FingerPrint,
			IsAvailable:  row.IsAvailable,
			CreationTime: row.CreationTime,
		})
	}

	return result, nil
}
