package product_removing

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UseCase struct {
	trManager     trManager
	receptionRepo receptionRepo
	productRepo   productRepo
	pvzLocker     pvzLocker
}

func New(trManager trManager, receptionRepo receptionRepo, pvzLocker pvzLocker, productRepo productRepo) (*UseCase, error) {
	if trManager == nil {
		return nil, errors.New("trManager is nil")
	}
	if receptionRepo == nil {
		return nil, errors.New("receptionRepo is nil")
	}
	if pvzLocker == nil {
		return nil, errors.New("pvzLocker is nil")
	}
	if productRepo == nil {
		return nil, errors.New("productRepo is nil")
	}
	return &UseCase{
		trManager:     trManager,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		pvzLocker:     pvzLocker,
	}, nil
}

func (uc *UseCase) RemoveLastProduct(ctx context.Context, pvzID model.PVZID) error {
	err := uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		err = uc.pvzLocker.Lock(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("pvzLocker.Lock: %w", err)
		}

		reception, err := uc.receptionRepo.GetInProgress(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("receptionRepo.GetInProgress: %w", err)
		}

		err = uc.productRepo.RemoveLast(ctx, reception.ID)
		if err != nil {
			return fmt.Errorf("productRepo.RemoveLast: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("trManager.Do: %w", err)
	}

	return nil
}
