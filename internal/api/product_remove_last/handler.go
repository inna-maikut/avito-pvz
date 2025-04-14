package product_remove_last

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type Handler struct {
	productRemoving productRemoving
	logger          internal.Logger
}

func New(productRemoving productRemoving, logger internal.Logger) (*Handler, error) {
	if productRemoving == nil {
		return nil, errors.New("productRemoving is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		productRemoving: productRemoving,
		logger:          logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleEmployee {
		api_handler.Forbidden(w, "only a user with the employee role can remove product")
		return
	}

	pvzID, err := model.ParsePVZID(r.PathValue("pvzId"))
	if err != nil {
		api_handler.BadRequest(w, "invalid pvzId")
		return
	}

	err = h.productRemoving.RemoveLastProduct(ctx, pvzID)
	if err != nil {
		err = fmt.Errorf("productRemoving.RemoveLastProduct(: %w", err)
		h.logger.Error("POST /pvz/{pvzId}/delete_last_product: internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("pvzId", pvzID))
		api_handler.InternalError(w, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
