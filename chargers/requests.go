package chargers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	CreateChargerRequest struct {
		Name     string   `json:"name"`
		Location Location `json:"location"`
	}
	CreateChargerResponse struct {
		Status string `json:"status"`
	}
	GetChargerRequest struct {
		Id string `json:"id"`
	}
	GetChargerResponse struct {
		Name          string              `json:"name"`
		Location      Location            `json:"location"`
		AverageRating float64             `json:"averageRating"`
		Ratings       []Rating            `json:"ratings"`
		Comments      []Comment           `json:"comments"`
		Reservations  []Reservation       `json:"reservations"`
		ClosestCity   string              `json:"closestCity"`
		Temperature   float64             `json:"temperature"`
		Weather       weatherAPIcondition `json:"weather"`
	}
	GetChargersRequest struct {
	}
	GetChargersResponse struct {
		Chargers []Charger `json:"chargers"`
	}
	UpdateChargerRequest struct {
		Id            string   `json:"id"`
		Name          string   `json:"name"`
		Location      Location `json:"location"`
		AverageRating float64  `json:"averageRating"`
	}
	UpdateChargerResponse struct {
		Status string `json:"status"`
	}
	DeleteChargerRequest struct {
		Id string `json:"id"`
	}
	DeleteChargerResponse struct {
		Status string `json:"status"`
	}
	//OTHER
	GetChargerRatingsRequest struct {
	}
	GetChargerRatingsResponse struct {
		Ratings []Rating `json:"ratings"`
	}
	GetChargerCommentsRequest struct {
	}
	GetChargerCommentsResponse struct {
		Comments []Comment `json:"comments"`
	}
	GetChargerReservationsRequest struct {
	}
	GetChargerReservationsResponse struct {
		Reservations []Reservation `json:"reservations"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateChargerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := CreateChargerRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeUpdateChargerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := UpdateChargerRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetChargerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetChargerRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetChargersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetChargerRequest{}
	return req, nil
}
func decodeDeleteChargerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := DeleteChargerRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
