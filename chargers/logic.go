package chargers

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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

func (s service) CreateCharger(ctx context.Context, name string) (string, error) {
	logger := log.With(s.logger, "method: ", "CreateCharger")
	charger := Charger{
		Name: name,
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
