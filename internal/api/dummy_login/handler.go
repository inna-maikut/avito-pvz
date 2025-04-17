package dummy_login

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type Handler struct {
	authenticating authenticating
	logger         internal.Logger
}

func New(authenticating authenticating, logger internal.Logger) (*Handler, error) {
	if authenticating == nil {
		return nil, errors.New("authenticating is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		authenticating: authenticating,
		logger:         logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var authRequest api.PostDummyLoginJSONBody
	if ok := api_handler.Parse(r, w, &authRequest); !ok {
		return
	}

	role, err := model.ParseUserRole(string(authRequest.Role))
	if err != nil {
		api_handler.BadRequest(w, "invalid role")
		return
	}

	token, err := h.authenticating.Auth(r.Context(), role)
	if err != nil {
		err = fmt.Errorf("authenticating.Auth: %w", err)
		h.logger.Error("POST /dummyLogin internal error", zap.Error(err), zap.Any("request", authRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}

	api_handler.OK(w, token)
}
