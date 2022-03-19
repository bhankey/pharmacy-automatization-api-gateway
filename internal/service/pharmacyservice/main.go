package pharmacyservice

import (
	"context"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

type Service struct {
	pharmacyStorage PharmacyStorage
}

type PharmacyStorage interface {
	CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error
	GetPharmacies(ctx context.Context, lastPharmacyID int, limit int) ([]entities.Pharmacy, error)
}

func NewPharmacyService(
	pharmacyRepo PharmacyStorage,
) *Service {
	return &Service{
		pharmacyStorage: pharmacyRepo,
	}
}
