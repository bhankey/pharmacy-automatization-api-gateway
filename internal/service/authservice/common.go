package authservice

import (
	"context"
	"time"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/golang-jwt/jwt/v4"
)

func (s *AuthService) createAndSaveRefreshToken(
	ctx context.Context,
	userID int,
	email string,
	role entities.Role,
	pharmacyID int,
	identifyData entities.UserIdentifyData,
) (string, error) {
	signedToken, err := s.createAndSignedToken(userID, email, role, pharmacyID, jwtExpireRefreshTime)
	if err != nil {
		return "", err
	}

	refreshToken := entities.RefreshToken{
		UserID:      userID,
		Token:       signedToken,
		UserAgent:   identifyData.UserAgent,
		IP:          identifyData.IP,
		FingerPrint: identifyData.FingerPrint,
	}

	if err := s.tokenStorage.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err // nolint: wrapcheck, nolintlint
	}

	return signedToken, nil
}

func (s *AuthService) createAndSignedToken(
	userID int,
	email string,
	role entities.Role,
	pharmacyID int,
	ttl time.Duration,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &entities.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email:      email,
		UserID:     userID,
		Role:       role,
		PharmacyID: pharmacyID,
	})

	signedToken, err := token.SignedString([]byte(s.jwtKey))
	if err != nil {
		return "", err // nolint: wrapcheck, nolintlint
	}

	return signedToken, nil
}
