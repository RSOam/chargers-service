package chargers

import (
	"context"
)

type ChargersService interface {
	CreateCharger(ctx context.Context, name string) (string, error)
	GetCharger(ctx context.Context, id string) (Charger, error)
}
