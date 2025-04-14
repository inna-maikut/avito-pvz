package pvz_list_getting

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

func TestUseCase_GetPVZList(t *testing.T) {
	type mocks struct {
		pvzRepo       *MockpvzRepo
		receptionRepo *MockreceptionRepo
		productRepo   *MockproductRepo
	}
	type args struct {
		receptedAtFrom, receptedAtTo *time.Time
		page, limit                  int64
	}

	from := time.Now().Add(-time.Hour * 24)
	to := time.Now()

	receptionID1 := model.NewReceptionID()
	receptionID2 := model.NewReceptionID()
	receptionID3 := model.NewReceptionID()

	pvzID1 := model.NewPVZID()
	pvzID2 := model.NewPVZID()

	productID1 := model.NewProductID()
	productID2 := model.NewProductID()

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes model.PVZList
		wantErr error
	}{
		{
			name: "success",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), &from, &to, int64(60), int64(30)).
					Return([]model.Reception{
						{
							ID:              receptionID1,
							PVZID:           pvzID1,
							ReceptionStatus: model.ReceptionStatusInProgress,
							ReceptedAt:      from,
						},
						{
							ID:              receptionID2,
							PVZID:           pvzID2,
							ReceptionStatus: model.ReceptionStatusClose,
							ReceptedAt:      to,
						},
						{
							ID:              receptionID3,
							PVZID:           pvzID2,
							ReceptionStatus: model.ReceptionStatusInProgress,
							ReceptedAt:      from,
						},
					}, nil)
				m.pvzRepo.EXPECT().
					Get(gomock.Any(), []model.PVZID{
						pvzID1, pvzID2,
					}).
					Return([]model.PVZ{
						{
							ID:           pvzID1,
							City:         "test1",
							RegisteredAt: from,
						},
						{
							ID:           pvzID2,
							City:         "test2",
							RegisteredAt: to,
						},
					}, nil)
				m.productRepo.EXPECT().
					GetByReceptionIDs(gomock.Any(), []model.ReceptionID{
						receptionID1, receptionID2, receptionID3,
					}).
					Return([]model.Product{
						{
							ID:          productID1,
							ReceptionID: receptionID1,
							Category:    model.ProductCategoryClothes,
							AddedAt:     from,
						},
						{
							ID:          productID2,
							ReceptionID: receptionID2,
							Category:    model.ProductCategoryElectronics,
							AddedAt:     to,
						},
					}, nil)
			},
			args: args{
				receptedAtFrom: &from,
				receptedAtTo:   &to,
				page:           3,
				limit:          30,
			},
			wantRes: model.PVZList{
				PVZs: []model.PVZ{
					{
						ID:           pvzID1,
						City:         "test1",
						RegisteredAt: from,
					},
					{
						ID:           pvzID2,
						City:         "test2",
						RegisteredAt: to,
					},
				},
				Receptions: []model.Reception{
					{
						ID:              receptionID1,
						PVZID:           pvzID1,
						ReceptionStatus: model.ReceptionStatusInProgress,
						ReceptedAt:      from,
					},
					{
						ID:              receptionID2,
						PVZID:           pvzID2,
						ReceptionStatus: model.ReceptionStatusClose,
						ReceptedAt:      to,
					},
					{
						ID:              receptionID3,
						PVZID:           pvzID2,
						ReceptionStatus: model.ReceptionStatusInProgress,
						ReceptedAt:      from,
					},
				},
				Products: []model.Product{
					{
						ID:          productID1,
						ReceptionID: receptionID1,
						Category:    model.ProductCategoryClothes,
						AddedAt:     from,
					},
					{
						ID:          productID2,
						ReceptionID: receptionID2,
						Category:    model.ProductCategoryElectronics,
						AddedAt:     to,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "empty_result",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), nil, nil, int64(30), int64(30)).
					Return([]model.Reception{}, nil)
			},
			args: args{
				receptedAtFrom: nil,
				receptedAtTo:   nil,
				page:           2,
				limit:          30,
			},
			wantRes: model.PVZList{
				PVZs:       nil,
				Receptions: nil,
				Products:   nil,
			},
			wantErr: nil,
		},
		{
			name: "error.Search",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), nil, nil, int64(30), int64(30)).
					Return(nil, assert.AnError)
			},
			args: args{
				receptedAtFrom: nil,
				receptedAtTo:   nil,
				page:           2,
				limit:          30,
			},
			wantRes: model.PVZList{},
			wantErr: assert.AnError,
		},
		{
			name: "error.GetPVZ",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), &from, &to, int64(60), int64(30)).
					Return([]model.Reception{
						{
							ID:              receptionID1,
							PVZID:           pvzID1,
							ReceptionStatus: model.ReceptionStatusInProgress,
							ReceptedAt:      from,
						},
					}, nil)
				m.pvzRepo.EXPECT().
					Get(gomock.Any(), []model.PVZID{
						pvzID1,
					}).
					Return(nil, assert.AnError)
				m.productRepo.EXPECT().
					GetByReceptionIDs(gomock.Any(), []model.ReceptionID{
						receptionID1,
					}).
					Return([]model.Product{}, nil)
			},
			args: args{
				receptedAtFrom: &from,
				receptedAtTo:   &to,
				page:           3,
				limit:          30,
			},
			wantRes: model.PVZList{},
			wantErr: assert.AnError,
		},
		{
			name: "error.GetPVZ",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), &from, &to, int64(60), int64(30)).
					Return([]model.Reception{
						{
							ID:              receptionID1,
							PVZID:           pvzID1,
							ReceptionStatus: model.ReceptionStatusInProgress,
							ReceptedAt:      from,
						},
					}, nil)
				m.pvzRepo.EXPECT().
					Get(gomock.Any(), []model.PVZID{
						pvzID1,
					}).
					Return(nil, assert.AnError)
				m.productRepo.EXPECT().
					GetByReceptionIDs(gomock.Any(), []model.ReceptionID{
						receptionID1,
					}).
					Return([]model.Product{}, nil)
			},
			args: args{
				receptedAtFrom: &from,
				receptedAtTo:   &to,
				page:           3,
				limit:          30,
			},
			wantRes: model.PVZList{},
			wantErr: assert.AnError,
		},
		{
			name: "error.GetByReceptionIDs",
			prepare: func(m *mocks) {
				m.receptionRepo.EXPECT().
					Search(gomock.Any(), &from, &to, int64(60), int64(30)).
					Return([]model.Reception{
						{
							ID:              receptionID1,
							PVZID:           pvzID1,
							ReceptionStatus: model.ReceptionStatusInProgress,
							ReceptedAt:      from,
						},
					}, nil)
				m.pvzRepo.EXPECT().
					Get(gomock.Any(), []model.PVZID{
						pvzID1,
					}).
					Return([]model.PVZ{}, nil)
				m.productRepo.EXPECT().
					GetByReceptionIDs(gomock.Any(), []model.ReceptionID{
						receptionID1,
					}).
					Return(nil, assert.AnError)
			},
			args: args{
				receptedAtFrom: &from,
				receptedAtTo:   &to,
				page:           3,
				limit:          30,
			},
			wantRes: model.PVZList{},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				pvzRepo:       NewMockpvzRepo(ctrl),
				receptionRepo: NewMockreceptionRepo(ctrl),
				productRepo:   NewMockproductRepo(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.pvzRepo, m.receptionRepo, m.productRepo)
			require.NoError(t, err)

			res, err := uc.GetPVZList(context.Background(), tc.args.receptedAtFrom, tc.args.receptedAtTo, tc.args.page, tc.args.limit)

			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}
