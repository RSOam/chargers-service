package chargers

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type service struct {
	db     ChargerDB
	logger log.Logger
}

func NewService(db ChargerDB, logger log.Logger) ChargersService {
	return &service{
		db:     db,
		logger: logger,
	}
}

func (s service) CreateCharger(ctx context.Context, name string, location Location) (string, error) {
	logger := log.With(s.logger, "method: ", "CreateCharger")
	charger := Charger{
		Name:          name,
		Location:      location,
		AverageRating: 0.0,
		Ratings:       []Rating{},
		Comments:      []Comment{},
		Reservations:  []Reservation{},
		Created:       time.Now().Format(time.RFC3339),
		Modified:      time.Now().Format(time.RFC3339),
	}
	if err := s.db.CreateCharger(ctx, charger); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("create Charger", nil)
	return "Ok", nil
}
func (s service) GetCharger(ctx context.Context, id string) (Charger, error) {
	logger := log.With(s.logger, "method", "GetCharger")
	charger, err := s.db.GetCharger(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return charger, err
	}
	logger.Log("Get Charger", id)
	return charger, nil
}
func (s service) GetChargers(ctx context.Context) ([]Charger, error) {
	logger := log.With(s.logger, "method", "GetChargers")
	chargers, err := s.db.GetChargers(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		return chargers, err
	}
	logger.Log("Get Chargers")
	return chargers, nil
}
func (s service) DeleteCharger(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteCharger")
	err := s.db.DeleteCharger(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("Delete Charger", id)
	return "Ok", nil
}
