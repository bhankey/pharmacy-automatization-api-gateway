package userhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bhankey/go-utils/pkg/apperror"
	deliveryhttp "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/models"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/go-openapi/strfmt"
)

func (h *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.RegisterRequest

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

	user := entities.User{
		ID:       0,
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email.String(),
		Password: *req.Password,

		Role:              entities.Role(*req.Role),
		DefaultPharmacyID: int(req.DefaultPharmacyID),
	}
	if err := h.userSrv.Registry(ctx, user); err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}

func (h *UserHandler) requestToChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.RequestPasswordChangeRequest

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

	if err := h.userSrv.RequestToResetPassword(ctx, req.Email.String()); err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}

func (h *UserHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.PasswordChangeRequest

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

	if err := h.userSrv.ResetPassword(ctx, req.Email.String(), *req.Code, *req.NewPassword); err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}

func (h *UserHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	lastUserID, err := strconv.Atoi(query.Get("last_id"))
	if err != nil || lastUserID < 0 {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 0 || limit > 5000 {
		h.WriteErrorResponse(ctx, w, apperror.NewClientError(apperror.WrongRequest, err))

		return
	}

	batch, err := h.userSrv.GetBatchOfUsers(ctx, lastUserID, limit)
	if err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	result := make(models.UserAllResponse, 0, len(batch))
	for _, user := range batch {
		result = append(result, &models.User{
			DefaultPharmacyID: int64(user.DefaultPharmacyID),
			Email:             strfmt.Email(user.Email),
			ID:                int64(user.ID),
			Name:              user.Name,
			Role:              string(user.Role),
			Surname:           user.Surname,
		})
	}

	deliveryhttp.WriteResponse(w, result)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.User

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

	if err := h.userSrv.UpdateUser(ctx, entities.User{
		DefaultPharmacyID: int(req.DefaultPharmacyID),
		Email:             string(req.Email),
		Name:              req.Name,
		Role:              entities.Role(req.Role),
		Surname:           req.Surname,
	}); err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}
