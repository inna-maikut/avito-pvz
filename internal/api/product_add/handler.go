package product_add

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal"
	"github.com/inna-maikut/avito-pvz/internal/api"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/model"
)

type Handler struct {
	productAdding productAdding
	logger        internal.Logger
}

func New(productAdding productAdding, logger internal.Logger) (*Handler, error) {
	if productAdding == nil {
		return nil, errors.New("productAdding is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		productAdding: productAdding,
		logger:        logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	if tokenInfo.UserRole != model.UserRoleEmployee {
		api_handler.Forbidden(w, "only a user with the employee role can add product")
		return
	}

	var request api.PostProductsJSONBody
	if ok := api_handler.Parse(r, w, &request); !ok {
		return
	}

	pvzID := model.PVZID(request.PvzId)
	category, err := model.ParseProductCategory(string(request.Type))
	if err != nil {
		api_handler.BadRequest(w, "invalid type")
		return
	}

	product, err := h.productAdding.AddProduct(ctx, pvzID, category)
	if err != nil {
		err = fmt.Errorf("productAdding.AddProduct: %w", err)
		h.logger.Error("POST /products/ internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("request", request))
		api_handler.InternalError(w, "internal server error")
		return
	}

	ID := product.ID.UUID()
	api_handler.Created(w, api.Product{
		Id:          &ID,
		ReceptionId: product.ReceptionID.UUID(),
		Type:        api.ProductType(product.Category.String()),
		DateTime:    &product.AddedAt,
	})
}
