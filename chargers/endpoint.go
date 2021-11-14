package chargers

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateCharger endpoint.Endpoint
	GetCharger    endpoint.Endpoint
}

func MakeEndpoints(s ChargersService) Endpoints {
	return Endpoints{
		CreateCharger: makeCreateChargerEndpoint(s),
		GetCharger:    makeGetChargerEndpoint(s),
	}
}

func makeCreateChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateChargerRequest)
		status, err := s.CreateCharger(ctx, req.Name)
		return CreateChargerResponse{Status: status}, err
	}
}
func makeGetChargerEndpoint(s ChargersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetChargerRequest)
		charger, err := s.GetCharger(ctx, req.Id)

		return GetChargerResponse{
			Name: charger.Name,
		}, err
	}
}
