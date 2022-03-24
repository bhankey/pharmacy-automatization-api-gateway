package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bhankey/go-utils/pkg/apperror"
	"github.com/bhankey/go-utils/pkg/logger"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/models"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/sirupsen/logrus"
)

type BaseHandler struct {
	Logger logger.Logger
}

func NewHandler(l logger.Logger) *BaseHandler {
	h := &BaseHandler{
		Logger: l,
	}

	return h
}

func (h *BaseHandler) WriteErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	h.Logger.WithFields(logrus.Fields{
		"error":      err,
		"request_id": ctx.Value(entities.RequestID),
	}).Errorf("response.error")

	var resp models.BaseResponse

	var clientError apperror.ClientError
	if errors.As(err, &clientError) {
		w.WriteHeader(clientError.GetHTTPCode())
		resp = models.BaseResponse{
			Error:   clientError.Message,
			Success: false,
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)

		resp = models.BaseResponse{
			Error:   "Something went wrong",
			Success: false,
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func WriteResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
