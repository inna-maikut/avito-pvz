package reception_creating

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UseCase struct {
	trManager     trManager
	receptionRepo receptionRepo
	pvzLocker     pvzLocker
	metric        metrics
}

func New(trManager trManager, receptionRepo receptionRepo, pvzLocker pvzLocker, metric metrics) (*UseCase, error) {
	if trManager == nil {
		return nil, errors.New("trManager is nil")
	}
	if receptionRepo == nil {
		return nil, errors.New("receptionRepo is nil")
	}
	if pvzLocker == nil {
		return nil, errors.New("pvzLocker is nil")
	}
	if metric == nil {
		return nil, errors.New("metric is nil")
	}

	return &UseCase{
		trManager:     trManager,
		receptionRepo: receptionRepo,
		pvzLocker:     pvzLocker,
		metric:        metric,
	}, nil
}

func (uc *UseCase) CreateReception(ctx context.Context, pvzID model.PVZID) (model.Reception, error) {
	var reception model.Reception

	err := uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		err = uc.pvzLocker.Lock(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("pvzLocker.Lock: %w", err)
		}

		_, err = uc.receptionRepo.GetInProgress(ctx, pvzID)
		if err == nil {
			return model.ErrReceptionAlreadyExists
		}
		if !errors.Is(err, model.ErrReceptionNotFound) {
			return fmt.Errorf("receptionRepo.GetInProgress: %w", err)
		}

		reception, err = uc.receptionRepo.Create(ctx, pvzID, model.ReceptionStatusInProgress)
		if err != nil {
			return fmt.Errorf("receptionRepo.Create: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Reception{}, fmt.Errorf("trManager.Do: %w", err)
	}

	uc.metric.ReceptionCreatedCountInc()

	return reception, nil
}
