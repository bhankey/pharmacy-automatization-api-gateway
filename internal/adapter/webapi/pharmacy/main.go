package pharmacy

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/api/pharmacyproto"
)

type APIClient struct {
	client pharmacyproto.PharmacyServiceClient
}

func NewPharmacyAPIClient(client pharmacyproto.PharmacyServiceClient) *APIClient {
	return &APIClient{
		client: client,
	}
}

func (c *APIClient) CreatePharmacy(ctx context.Context, pharmacy entities.Pharmacy) error {
	errBase := fmt.Sprintf("pharmacy.CreatePharmacy( %v)", pharmacy)

	_, err := c.client.CreatePharmacy(ctx, &pharmacyproto.NewPharmacy{
		Name: pharmacy.Name,
		Address: &pharmacyproto.Address{
			City:   pharmacy.Address.City,
			Street: pharmacy.Address.Street,
			House:  pharmacy.Address.House,
		},
	})
	if err != nil {
		return fmt.Errorf("%s, failed to create pharmacy: %w", errBase, err)
	}

	return nil
}

func (c *APIClient) GetPharmacies(ctx context.Context, lastID, limit int) ([]entities.Pharmacy, error) {
	errBase := fmt.Sprintf("pharmacy.GetPharmacies(%d, %d)", lastID, limit)

	pharmacies, err := c.client.GetPharmacies(ctx, &pharmacyproto.PaginationRequest{
		LastId: int64(lastID),
		Limit:  int64(limit),
	},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get pharmacies: %w", errBase, err)
	}

	res := make([]entities.Pharmacy, 0, len(pharmacies.Pharmacies))
	for _, pharmacy := range pharmacies.Pharmacies {
		res = append(res, entities.Pharmacy{
			ID:        int(pharmacy.GetId()),
			Name:      pharmacy.GetName(),
			IsBlocked: false, // TODO
			Address: entities.Address{
				City:   pharmacy.GetAddress().GetCity(),
				Street: pharmacy.GetAddress().GetStreet(),
				House:  pharmacy.GetAddress().GetHouse(),
			},
		})
	}

	return res, nil
}
