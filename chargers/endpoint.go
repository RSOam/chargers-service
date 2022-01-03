package chargers

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateCharger endpoint.Endpoint
	GetCharger    endpoint.Endpoint
	GetChargers   endpoint.Endpoint
	UpdateCharger endpoint.Endpoint
	DeleteCharger endpoint.Endpoint
}

func MakeEndpoints(s ChargersService) Endpoints {
	return Endpoints{
		CreateCharger: makeCreateChargerEndpoint(s),
		GetCharger:    makeGetChargerEndpoint(s),
		GetChargers:   makeGetChargersEndpoint(s),
		UpdateCharger: makeUpdateChargerEndpoint(s),
		DeleteCharger: makeDeleteChargerEndpoint(s),
	}
}

func makeCreateChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateChargerRequest)
		status, err := s.CreateCharger(ctx, req.Name, req.Location)
		return CreateChargerResponse{Status: status}, err
	}
}
func makeGetChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetChargerRequest)
		charger, err := s.GetCharger(ctx, req.Id)
		return GetChargerResponse{
			Name:          charger.Name,
			Location:      charger.Location,
			AverageRating: charger.AverageRating,
			Ratings:       charger.Ratings,
			Comments:      charger.Comments,
			Reservations:  charger.Reservations,
		}, err
	}
}
func makeGetChargersEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		chargers, err := s.GetChargers(ctx)
		return GetChargersResponse{
			Chargers: chargers,
		}, err
	}
}
func makeDeleteChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteChargerRequest)
		status, err := s.DeleteCharger(ctx, req.Id)
		return DeleteChargerResponse{
			Status: status,
		}, err
	}
}
func makeUpdateChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateChargerRequest)
		status, err := s.UpdateCharger(ctx, req.Id, req.Name, req.Location, req.AverageRating)
		return UpdateChargerResponse{Status: status}, err
	}
}
