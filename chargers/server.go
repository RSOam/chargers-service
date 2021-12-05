package chargers

import (
	"context"
	"net/http"

	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHttpServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/chargers").Handler(ht.NewServer(
		endpoints.CreateCharger,
		decodeCreateChargerRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/chargers/{id}").Handler(ht.NewServer(
		endpoints.GetCharger,
		decodeGetChargerRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/chargers").Handler(ht.NewServer(
		endpoints.GetChargers,
		decodeGetChargersRequest,
		encodeResponse,
	))
	r.Methods("DELETE").Path("/chargers/{id}").Handler(ht.NewServer(
		endpoints.DeleteCharger,
		decodeDeleteChargerRequest,
		encodeResponse,
	))
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
