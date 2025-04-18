package pvz_registering

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UseCase struct {
	pvzRepo pvzRepo
	metric  metrics
}

func New(pvzRepo pvzRepo, metric metrics) (*UseCase, error) {
	if pvzRepo == nil {
		return nil, errors.New("pvzRepo is nil")
	}
	if metric == nil {
		return nil, errors.New("metric is nil")
	}

	return &UseCase{
		pvzRepo: pvzRepo,
		metric:  metric,
	}, nil
}

func (uc *UseCase) RegisterPVZ(ctx context.Context, city string) (model.PVZ, error) {
	pvz := model.PVZ{
		ID:           model.NewPVZID(),
		City:         city,
		RegisteredAt: time.Now(),
	}

	err := uc.pvzRepo.Register(ctx, pvz)
	if err != nil {
		return model.PVZ{}, fmt.Errorf("pvzRepo.Register: %w", err)
	}

	uc.metric.PVZRegisteredCountInc()

	return pvz, nil
}
