package pharmacyservice

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

func (s *Service) CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	errBase := fmt.Sprintf("pharmacyservice.CreatePharmacy(%v)", pharmacy)

	if err := s.pharmacyStorage.CreatePharmacy(ctx, pharmacy); err != nil {
		return fmt.Errorf("%s: failed to create pharmacy: %w", errBase, err)
	}

	return nil
}

func (s *Service) GetBatchOfPharmacies(
	ctx context.Context,
	lastPharmacyID int,
	limit int,
) ([]entities.Pharmacy, error) {
	errBase := fmt.Sprintf("pharmacyservice.GetBatchOfPharmacies(%d, %d)", lastPharmacyID, limit)

	pharmacies, err := s.pharmacyStorage.GetPharmacies(ctx, lastPharmacyID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get batch of pharmacies: %w", errBase, err)
	}

	return pharmacies, nil
}
