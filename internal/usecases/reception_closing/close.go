package reception_closing

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
}

func New(trManager trManager, receptionRepo receptionRepo, pvzLocker pvzLocker) (*UseCase, error) {
	if trManager == nil {
		return nil, errors.New("trManager is nil")
	}
	if receptionRepo == nil {
		return nil, errors.New("receptionRepo is nil")
	}
	if pvzLocker == nil {
		return nil, errors.New("pvzLocker is nil")
	}

	return &UseCase{
		trManager:     trManager,
		receptionRepo: receptionRepo,
		pvzLocker:     pvzLocker,
	}, nil
}

func (uc *UseCase) CloseReception(ctx context.Context, pvzID model.PVZID) (model.Reception, error) {
	var reception model.Reception

	err := uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		err = uc.pvzLocker.Lock(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("pvzLocker.Lock: %w", err)
		}

		reception, err = uc.receptionRepo.GetInProgress(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("receptionRepo.GetInProgress: %w", err)
		}

		reception.ReceptionStatus = model.ReceptionStatusClose

		err = uc.receptionRepo.SetStatus(ctx, reception.ID, model.ReceptionStatusClose)
		if err != nil {
			return fmt.Errorf("receptionRepo.SetStatus: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Reception{}, fmt.Errorf("trManager.Do: %w", err)
	}

	return reception, nil
}
