package chargers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/heptiolabs/healthcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewHttpServer(ctx context.Context, endpoints Endpoints, db *mongo.Database, consul consulapi.Client) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	r.Use(cors)
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("database", DatabasePingCheck(db, 1*time.Second))
	//health.AddLivenessCheck("sick", simulateSick(consul))
	r.Methods("POST").Path("/chargers").Handler(ht.NewServer(
		endpoints.CreateCharger,
		decodeCreateChargerRequest,
		encodeResponse,
	))
	r.Methods("PUT").Path("/chargers/{id}").Handler(ht.NewServer(
		endpoints.UpdateCharger,
		decodeUpdateChargerRequest,
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
	r.HandleFunc("/health/live", health.LiveEndpoint)
	r.HandleFunc("/health/ready", health.ReadyEndpoint)
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func DatabasePingCheck(database *mongo.Database, timeout time.Duration) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if database == nil {
			return fmt.Errorf("database is nil")
		}
		return database.Client().Ping(ctx, readpref.Primary())
	}
}
func simulateSick(consul consulapi.Client) healthcheck.Check {
	return func() error {
		str, err := getConsulCheck(consul, "health")
		if err != nil {
			return err
		}
		if str != "true" {
			return fmt.Errorf("service sick")
		} else {
			return nil
		}
	}
}
func getConsulCheck(consul consulapi.Client, key string) (string, error) {
	kv := consul.KV()
	keyPair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}
	return string(keyPair.Value), nil
}
