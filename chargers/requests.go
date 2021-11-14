package chargers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	CreateChargerRequest struct {
		Name string `json:"name"`
	}
	CreateChargerResponse struct {
		Status string `json:"status"`
	}
	GetChargerRequest struct {
		Id string `json:"id"`
	}
	GetChargerResponse struct {
		Name string `json:"name"`
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
func decodeGetChargerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetChargerRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
