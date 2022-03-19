package pharmacyhandler

import (
	"encoding/json"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/apperror"
	deliveryhttp "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/models"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"
)

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.Pharmacy

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&req)
	if err != nil {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	if err := req.Validate(strfmt.NewFormats()); err != nil {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	pharmacy := entities.Pharmacy{
		Name:      req.Name,
		IsBlocked: false,
		Address: entities.Address{
			City:   req.Address.City,
			Street: req.Address.Street,
			House:  req.Address.House,
		},
	}
	if err := h.pharmacySrv.CreatePharmacy(ctx, pharmacy); err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}

func (h *Handler) all(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	lastID, err := strconv.Atoi(query.Get("last_id"))
	if err != nil || lastID < 0 {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 0 || limit > 5000 {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	batch, err := h.pharmacySrv.GetBatchOfPharmacies(ctx, lastID, limit)
	if err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	resp := make([]*models.Pharmacy, 0, len(batch))

	for _, pharmacy := range batch {
		resp = append(resp, &models.Pharmacy{
			Address: &models.Address{
				City:   pharmacy.Address.City,
				House:  pharmacy.Address.House,
				Street: pharmacy.Address.Street,
			},
			ID:   int64(pharmacy.ID),
			Name: pharmacy.Name,
		})
	}
	deliveryhttp.WriteResponse(w, models.PharmacyGetAllResponse{
		Pharmacys: resp,
	})
}
