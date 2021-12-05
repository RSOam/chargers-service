package chargers

import (
	"context"
)

type ChargersService interface {
	CreateCharger(ctx context.Context, name string, location Location) (string, error)
	GetCharger(ctx context.Context, id string) (Charger, error)
	GetChargers(ctx context.Context) ([]Charger, error)
	DeleteCharger(ctx context.Context, id string) (string, error)
}
