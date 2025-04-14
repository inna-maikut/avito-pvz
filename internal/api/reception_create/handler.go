package reception_create

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
	receptionCreating receptionCreating
	logger            internal.Logger
}

func New(receptionCreating receptionCreating, logger internal.Logger) (*Handler, error) {
	if receptionCreating == nil {
		return nil, errors.New("receptionCreating is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		receptionCreating: receptionCreating,
		logger:            logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleEmployee {
		api_handler.Forbidden(w, "only a user with the employee role can create reception")
		return
	}

	var createReceptionRequest api.PostReceptionsJSONBody
	if ok := api_handler.Parse(r, w, &createReceptionRequest); !ok {
		return
	}

	pvzID := model.PVZID(createReceptionRequest.PvzId)

	reception, err := h.receptionCreating.CreateReception(ctx, pvzID)
	if err != nil {
		err = fmt.Errorf("receptionCreating.CreateReception: %w", err)
		h.logger.Error("POST /receptions/ internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("request", createReceptionRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}

	ID := reception.ID.UUID()
	api_handler.Created(w, api.Reception{
		PvzId:    types.UUID(pvzID),
		Id:       &ID,
		Status:   api.ReceptionStatus(reception.ReceptionStatus.String()),
		DateTime: reception.ReceptedAt,
	})
}
