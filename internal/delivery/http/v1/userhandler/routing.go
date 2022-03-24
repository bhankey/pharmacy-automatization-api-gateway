package userhandler

import (
	"context"

	deliveryhttp "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/middleware"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Router  chi.Router
	userSrv UserSrv

	*deliveryhttp.BaseHandler
}

type UserSrv interface {
	Registry(ctx context.Context, user entities.User) error
	ResetPassword(ctx context.Context, email, code, newPassword string) error
	RequestToResetPassword(ctx context.Context, email string) error
	UpdateUser(ctx context.Context, user entities.User) error
	GetBatchOfUsers(ctx context.Context, lastClientID int, limit int) ([]entities.User, error)
}

func NewUserHandler(
	baseHandler *deliveryhttp.BaseHandler,
	userSrv UserSrv,
	authMiddleware *middleware.AuthMiddleware,
) *UserHandler {
	router := chi.NewRouter()

	handler := &UserHandler{
		Router:      router,
		userSrv:     userSrv,
		BaseHandler: baseHandler,
	}

	handler.initRoutes(router, authMiddleware)

	return handler
}

func (h *UserHandler) initRoutes(router chi.Router, authMiddleware *middleware.AuthMiddleware) {
	router.Use(authMiddleware.Middleware)

	router.Post("/register", h.register)
	router.Post("/request_to_change_password", h.requestToChangePassword)
	router.Post("/change_password", h.changePassword)
	router.Get("/all", h.getAll)
	router.Get("/update", h.update)
}
