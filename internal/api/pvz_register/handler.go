package pvz_register

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

var validCities = []string{"Москва", "Санкт-Петербург", "Казань"}

type Handler struct {
	pvzRegistering pvzRegistering
	logger         internal.Logger
}

func New(pvzRegistering pvzRegistering, logger internal.Logger) (*Handler, error) {
	if pvzRegistering == nil {
		return nil, errors.New("pvzRegistering is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		pvzRegistering: pvzRegistering,
		logger:         logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleModerator {
		api_handler.Forbidden(w, "only a user with the moderator role can create a pickup point in the system")
		return
	}

	var registerPVZRequest api.PVZ
	if ok := api_handler.Parse(r, w, &registerPVZRequest); !ok {
		return
	}

	city := string(registerPVZRequest.City)
	if !validateCity(city) {
		api_handler.BadRequest(w, "invalid city")
		return
	}

	pvz, err := h.pvzRegistering.RegisterPVZ(ctx, city)
	if err != nil {
		err = fmt.Errorf("pvzRegistering.RegisterPVZ: %w", err)
		h.logger.Error("POST /api/pvz internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("request", registerPVZRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}

	ID := pvz.ID.UUID()
	api_handler.Created(w, api.PVZ{
		Id:               &ID,
		City:             api.PVZCity(pvz.City),
		RegistrationDate: &pvz.RegisteredAt,
	})
}

func validateCity(city string) bool {
	return slices.Contains(validCities, city)
}
