package pvz_get

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

const (
	defaultLimit = 10
	maxLimit     = 30
)

type PVZGetResponsePVZItem struct {
	PVZ        api.PVZ                       `json:"pvz"`
	Receptions []PVZGetResponseReceptionItem `json:"receptions"`
}

type PVZGetResponseReceptionItem struct {
	Reception api.Reception `json:"reception"`
	Products  []api.Product `json:"products"`
}

type Handler struct {
	pvzListGetting pvzListGetting
	logger         internal.Logger
}

func New(pvzListGetting pvzListGetting, logger internal.Logger) (*Handler, error) {
	if pvzListGetting == nil {
		return nil, errors.New("pvzListGetting is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		pvzListGetting: pvzListGetting,
		logger:         logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleModerator && tokenInfo.UserRole != model.UserRoleEmployee {
		api_handler.Forbidden(w, "only a user with the moderator or employee role can create a pickup point in the system")
		return
	}

	from, to, page, limit, err := parseQuery(r.URL.Query())
	if err != nil {
		api_handler.BadRequest(w, "validation query: "+err.Error())
		return
	}

	pvzList, err := h.pvzListGetting.GetPVZList(ctx, from, to, page, limit)
	if err != nil {
		err = fmt.Errorf("pvzListGetting.RegisterPVZ: %w", err)
		h.logger.Error("GET /pvz internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("query", r.URL.Query()))
		api_handler.InternalError(w, "internal server error")
		return
	}

	api_handler.OK(w, convertToDTO(pvzList))
}

func ptrOf[T any](value T) *T {
	return &value
}

func parseQuery(query url.Values) (from, to *time.Time, page, limit int64, err error) {
	startDate := query.Get("startDate")
	if startDate != "" {
		var fromValue strfmt.DateTime
		fromValue, err = strfmt.ParseDateTime(startDate)
		if err != nil {
			return nil, nil, 0, 0, fmt.Errorf("parse start date: %w", err)
		}
		from = ptrOf(time.Time(fromValue))
	}

	endDate := query.Get("endDate")
	if endDate != "" {
		var toValue strfmt.DateTime
		toValue, err = strfmt.ParseDateTime(endDate)
		if err != nil {
			return nil, nil, 0, 0, fmt.Errorf("parse end date: %w", err)
		}
		to = ptrOf(time.Time(toValue))
	}

	pageParam := query.Get("page")
	if pageParam != "" {
		page, err = strconv.ParseInt(pageParam, 10, 64)
		if err != nil {
			return nil, nil, 0, 0, fmt.Errorf("parse page: %w", err)
		}
		if page < 1 {
			return nil, nil, 0, 0, fmt.Errorf("page must be greater than zero")
		}
	} else {
		page = 1
	}

	limitParam := query.Get("limit")
	if limitParam != "" {
		limit, err = strconv.ParseInt(limitParam, 10, 64)
		if err != nil {
			return nil, nil, 0, 0, fmt.Errorf("parse page: %w", err)
		}
		if limit < 1 {
			return nil, nil, 0, 0, fmt.Errorf("limit must be greater than zero")
		}
		if limit > maxLimit {
			return nil, nil, 0, 0, fmt.Errorf("limit must be not greater than 30")
		}
	} else {
		limit = defaultLimit
	}

	return from, to, page, limit, nil
}

func convertToDTO(pvzList model.PVZList) []PVZGetResponsePVZItem {
	res := make([]PVZGetResponsePVZItem, 0, len(pvzList.PVZs))

	productsByReception := make(map[model.ReceptionID][]api.Product, len(pvzList.PVZs))
	for _, product := range pvzList.Products {
		productsByReception[product.ReceptionID] = append(productsByReception[product.ReceptionID], api.Product{
			Id:          (*types.UUID)(&product.ID),
			ReceptionId: types.UUID(product.ReceptionID),
			Type:        api.ProductType(product.Category.String()),
			DateTime:    ptrOf(product.AddedAt),
		})
	}

	receptionsByPVZ := make(map[model.PVZID][]PVZGetResponseReceptionItem, len(pvzList.PVZs))
	for _, reception := range pvzList.Receptions {
		receptionsByPVZ[reception.PVZID] = append(receptionsByPVZ[reception.PVZID], PVZGetResponseReceptionItem{
			Reception: api.Reception{
				Id:       (*types.UUID)(&reception.ID),
				PvzId:    types.UUID(reception.PVZID),
				Status:   api.ReceptionStatus(reception.ReceptionStatus.String()),
				DateTime: reception.ReceptedAt,
			},
			Products: productsByReception[reception.ID],
		})
	}

	for _, pvz := range pvzList.PVZs {
		res = append(res, PVZGetResponsePVZItem{
			PVZ: api.PVZ{
				Id:               (*types.UUID)(&pvz.ID),
				City:             api.PVZCity(pvz.City),
				RegistrationDate: &pvz.RegisteredAt,
			},
			Receptions: receptionsByPVZ[pvz.ID],
		})
	}

	return res
}
