package pvz_list_getting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UseCase struct {
	pvzRepo       pvzRepo
	receptionRepo receptionRepo
	productRepo   productRepo
}

func New(
	pvzRepo pvzRepo,
	receptionRepo receptionRepo,
	productRepo productRepo,
) (*UseCase, error) {
	if pvzRepo == nil {
		return nil, errors.New("pvzRepo is nil")
	}
	if receptionRepo == nil {
		return nil, errors.New("receptionRepo is nil")
	}
	if productRepo == nil {
		return nil, errors.New("productRepo is nil")
	}
	return &UseCase{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}, nil
}

func (uc *UseCase) GetPVZList(ctx context.Context, receptedAtFrom, receptedAtTo *time.Time, page, limit int64) (model.PVZList, error) {
	offset := (page - 1) * limit
	receptions, err := uc.receptionRepo.Search(ctx, receptedAtFrom, receptedAtTo, offset, limit)
	if err != nil {
		return model.PVZList{}, fmt.Errorf("receptionRepo.Search: %w", err)
	}
	if len(receptions) == 0 {
		return model.PVZList{}, nil
	}

	receptionIDs := make([]model.ReceptionID, 0, len(receptions))
	pvzIDs := make([]model.PVZID, 0, len(receptions))
	pvzIDMap := make(map[model.PVZID]struct{}, len(receptions))
	for _, reception := range receptions {
		receptionIDs = append(receptionIDs, reception.ID)
		if _, found := pvzIDMap[reception.PVZID]; !found {
			pvzIDs = append(pvzIDs, reception.PVZID)
			pvzIDMap[reception.PVZID] = struct{}{}
		}
	}

	var (
		eg       *errgroup.Group
		PVZs     []model.PVZ
		products []model.Product
	)
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		PVZs, err = uc.pvzRepo.Get(ctx, pvzIDs)
		if err != nil {
			return fmt.Errorf("pvzRepo.Get: %w", err)
		}

		return nil
	})

	eg.Go(func() (err error) {
		products, err = uc.productRepo.GetByReceptionIDs(ctx, receptionIDs)
		if err != nil {
			return fmt.Errorf("productRepo.GetByReceptionIDs: %w", err)
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return model.PVZList{}, fmt.Errorf("errgroup.Wait: %w", err)
	}

	return model.PVZList{
		PVZs:       PVZs,
		Receptions: receptions,
		Products:   products,
	}, nil
}
