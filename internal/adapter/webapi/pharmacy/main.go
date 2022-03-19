package pharmacy

import (
	"context"

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
	_, err := c.client.CreatePharmacy(ctx, &pharmacyproto.NewPharmacy{
		Name: pharmacy.Name,
		Address: &pharmacyproto.Address{
			City:   pharmacy.Address.City,
			Street: pharmacy.Address.Street,
			House:  pharmacy.Address.House,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *APIClient) GetPharmacies(ctx context.Context, lastID, limit int) ([]entities.Pharmacy, error) {
	pharmacies, err := c.client.GetPharmacies(ctx, &pharmacyproto.PaginationRequest{
		LastId: int64(lastID),
		Limit:  int64(limit),
	},
	)
	if err != nil {
		return nil, err
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
