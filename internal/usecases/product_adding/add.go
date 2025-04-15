package product_adding

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
	metric        metrics
}

func New(trManager trManager, receptionRepo receptionRepo, pvzLocker pvzLocker, productRepo productRepo, metric metrics) (*UseCase, error) {
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
	if metric == nil {
		return nil, errors.New("metric is nil")
	}
	return &UseCase{
		trManager:     trManager,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		pvzLocker:     pvzLocker,
		metric:        metric,
	}, nil
}

func (uc *UseCase) AddProduct(ctx context.Context, pvzID model.PVZID, category model.ProductCategory) (model.Product, error) {
	var product model.Product

	err := uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		err = uc.pvzLocker.Lock(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("pvzLocker.Lock: %w", err)
		}

		reception, err := uc.receptionRepo.GetInProgress(ctx, pvzID)
		if err != nil {
			return fmt.Errorf("receptionRepo.GetInProgress: %w", err)
		}

		product, err = uc.productRepo.Create(ctx, reception.ID, category)
		if err != nil {
			return fmt.Errorf("productRepo.Create: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Product{}, fmt.Errorf("trManager.Do: %w", err)
	}

	uc.metric.ProductAddedCountInc()

	return product, nil
}
