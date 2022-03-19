package v1

import (
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/authhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/pharmacyhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/swaggerhandler"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1/userhandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(
	swaggerHandler *swaggerhandler.SwaggerHandler,
	authHandler *authhandler.AuthHandler,
	userHandler *userhandler.UserHandler,
	pharmacyHandler *pharmacyhandler.Handler,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))

	router.Mount("/docs", swaggerHandler.Router)
	router.Mount("/auth", authHandler.Router)
	router.Mount("/user", userHandler.Router)
	router.Mount("/pharmacy", pharmacyHandler.Router)

	return router
}
