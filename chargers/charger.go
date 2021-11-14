package chargers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Charger struct {
	ID   primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Name string             `json:"name"`
}

type ChargerDB interface {
	CreateCharger(ctx context.Context, charger Charger) error
	GetCharger(ctx context.Context, id string) (Charger, error)
}
