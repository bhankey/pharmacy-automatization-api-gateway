package middleware

import (
	"context"
	"net/http"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/entities"
)

const FingerPrintHeader = "X-Device-FingerPrint"

func FingerPrint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fingerPrint := r.Header.Get(FingerPrintHeader)

		ctx := r.Context()
		ctx = context.WithValue(ctx, entities.DeviceFingerPrint, fingerPrint)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
