package reception_close

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type Handler struct {
	receptionClosing receptionClosing
	logger           internal.Logger
}

func New(receptionClosing receptionClosing, logger internal.Logger) (*Handler, error) {
	if receptionClosing == nil {
		return nil, errors.New("receptionClosing is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		receptionClosing: receptionClosing,
		logger:           logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleEmployee {
		api_handler.Forbidden(w, "only a user with the employee role can close reception")
		return
	}

	pvzID, err := model.ParsePVZID(r.PathValue("pvzId"))
	if err != nil {
		api_handler.BadRequest(w, "invalid pvzId")
		return
	}

	reception, err := h.receptionClosing.CloseReception(ctx, pvzID)
	if err != nil {
		err = fmt.Errorf("receptionClosing.CloseReception: %w", err)
		h.logger.Error("POST /pvz/{pvzId}/close_last_reception: internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("pvzId", pvzID))
		api_handler.InternalError(w, "internal server error")
		return
	}

	ID := reception.ID.UUID()
	api_handler.OK(w, api.Reception{
		PvzId:    types.UUID(pvzID),
		Id:       &ID,
		Status:   api.ReceptionStatus(reception.ReceptionStatus.String()),
		DateTime: reception.ReceptedAt,
	})
}
