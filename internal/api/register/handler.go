package register

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/model"
	"github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	registering registering
	logger      internal.Logger
}

func New(registering registering, logger internal.Logger) (*Handler, error) {
	if registering == nil {
		return nil, errors.New("registering is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		registering: registering,
		logger:      logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var authRequest api.PostRegisterJSONBody
	if ok := api_handler.Parse(r, w, &authRequest); !ok {
		return
	}

	role, err := model.ParseUserRole(string(authRequest.Role))
	if err != nil {
		api_handler.BadRequest(w, "invalid role")
		return
	}

	if authRequest.Password == "" {
		api_handler.BadRequest(w, "empty password")
		return
	}

	user, err := h.registering.Register(r.Context(), string(authRequest.Email), authRequest.Password, role)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			api_handler.BadRequest(w, "user already exists")
			return
		}
		err = fmt.Errorf("registering.Register: %w", err)
		h.logger.Error("POST /register internal error", zap.Error(err), zap.Any("request", authRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}
	ID := user.UserID.UUID()
	api_handler.Created(w, api.User{
		Id:    &ID,
		Email: types.Email(user.Email),
		Role:  api.UserRole(model.UserRole.String(user.UserRole)),
	})
}
