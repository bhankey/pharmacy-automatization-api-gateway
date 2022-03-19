package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/apperror"
	deliveryhttp "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/models"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/go-openapi/strfmt"
)

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.UserLoginRequest

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

	deviceFingerPrint, _ := ctx.Value(entities.DeviceFingerPrint).(string)

	identifyData := entities.UserIdentifyData{
		IP:          r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		FingerPrint: deviceFingerPrint,
	}

	tokens, err := h.authSrv.Login(ctx, req.Email.String(), *req.Password, identifyData)
	if err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    tokens.AccessToken,
		Secure:   true,
		HttpOnly: true,
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    tokens.RefreshToken,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() { _ = r.Body.Close() }()
	var req models.UserRefreshRequest

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

	deviceFingerPrint, _ := ctx.Value(entities.DeviceFingerPrint).(string)

	identifyData := entities.UserIdentifyData{
		IP:          r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		FingerPrint: deviceFingerPrint,
	}

	tokens, err := h.authSrv.RefreshToken(ctx, *req.Token, identifyData)
	if err != nil {
		h.WriteErrorResponse(ctx, w, err)

		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    tokens.AccessToken,
		Secure:   true,
		HttpOnly: true,
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    tokens.RefreshToken,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	deliveryhttp.WriteResponse(w, models.BaseResponse{
		Error:   "",
		Success: true,
	})
}
