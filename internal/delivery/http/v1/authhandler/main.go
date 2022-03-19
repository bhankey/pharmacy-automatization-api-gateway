package authhandler

import (
	"context"

	deliveryhttp "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	Router  chi.Router
	authSrv AuthSrv

	*deliveryhttp.BaseHandler
}

type AuthSrv interface {
	Login(ctx context.Context, email, password string, identifyData entities.UserIdentifyData) (entities.Tokens, error)
	RefreshToken(ctx context.Context, refreshToken string, identifyData entities.UserIdentifyData) (entities.Tokens, error)
}

func NewAuthHandler(baseHandler *deliveryhttp.BaseHandler, authSrv AuthSrv) *AuthHandler {
	router := chi.NewRouter()

	handler := &AuthHandler{
		Router:      router,
		authSrv:     authSrv,
		BaseHandler: baseHandler,
	}

	handler.initRoutes(router)

	return handler
}

func (h *AuthHandler) initRoutes(router chi.Router) {
	router.Post("/login", h.login)
	router.Post("/refresh", h.refresh)
}
