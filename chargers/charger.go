package chargers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Charger struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Location      Location           `json:"location"`
	AverageRating float64            `json:"averageRating"`
	Ratings       []Rating           `json:"ratings"`
	Comments      []Comment          `json:"comments"`
	Reservations  []Reservation      `json:"reservations"`
	Created       string             `json:"created"`
	Modified      string             `json:"modified"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Text      string             `json:"text"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}
type Rating struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Rating    int                `json:"rating"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type Reservation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	From      string             `json:"from"`
	To        string             `json:"to"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type ChargerExtra struct {
	City    string              `json:"closestCity"`
	Temp    float64             `json:"temp"`
	Weather weatherAPIcondition `json:"weather"`
}
type ChargerDB interface {
	CreateCharger(ctx context.Context, charger Charger) error
	GetCharger(ctx context.Context, id string) (Charger, ChargerExtra, error)
	GetChargers(ctx context.Context) ([]Charger, error)
	UpdateCharger(ctx context.Context, id string, charger Charger) error
	DeleteCharger(ctx context.Context, id string) error
}
