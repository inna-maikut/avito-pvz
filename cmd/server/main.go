package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/api/dummy_login"
	"github.com/inna-maikut/avito-pvz/internal/api/product_add"
	"github.com/inna-maikut/avito-pvz/internal/api/product_remove_last"
	"github.com/inna-maikut/avito-pvz/internal/api/pvz_get"
	"github.com/inna-maikut/avito-pvz/internal/api/pvz_register"
	"github.com/inna-maikut/avito-pvz/internal/api/reception_close"
	"github.com/inna-maikut/avito-pvz/internal/api/reception_create"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/config"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/middleware"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/pg"
	"github.com/inna-maikut/avito-pvz/internal/repository"
	"github.com/inna-maikut/avito-pvz/internal/usecases/dummy_authenticating"
	"github.com/inna-maikut/avito-pvz/internal/usecases/product_adding"
	"github.com/inna-maikut/avito-pvz/internal/usecases/product_removing"
	"github.com/inna-maikut/avito-pvz/internal/usecases/pvz_list_getting"
	"github.com/inna-maikut/avito-pvz/internal/usecases/pvz_registering"
	"github.com/inna-maikut/avito-pvz/internal/usecases/reception_closing"
	"github.com/inna-maikut/avito-pvz/internal/usecases/reception_creating"
)

const (
	readHeaderTimeout = time.Second
)

func main() { //nolint:gocognit
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			if typedErr, ok := panicErr.(error); ok {
				logger.Error("panic error", zap.Error(typedErr))
			} else {
				logger.Error("panic", zap.Any("message", panicErr))
			}
		}

		_ = logger.Sync()
	}()

	db, cancelDB, err := pg.NewDB(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("unable to init database: %w", err))
	}
	defer cancelDB()

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))

	tokenProvider, err := jwt.NewProviderFromEnv()
	if err != nil {
		panic(fmt.Errorf("create jwt provider: %w", err))
	}

	productRepo, err := repository.NewProductRepository(db, trmsqlx.DefaultCtxGetter)
	if err != nil {
		panic(fmt.Errorf("create product repository: %w", err))
	}

	pvzRepo, err := repository.NewPVZRepository(db, trmsqlx.DefaultCtxGetter)
	if err != nil {
		panic(fmt.Errorf("create pvz repository: %w", err))
	}

	receptionRepo, err := repository.NewReceptionRepository(db, trmsqlx.DefaultCtxGetter)
	if err != nil {
		panic(fmt.Errorf("create reception repository: %w", err))
	}

	pvzLocker, err := repository.NewPVZLocker(db, trmsqlx.DefaultCtxGetter)
	if err != nil {
		panic(fmt.Errorf("create pvz locker: %w", err))
	}

	// Use cases

	dummyAuthentication, err := dummy_authenticating.New(tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create dummy_authenticating use case: %w", err))
	}

	productAdding, err := product_adding.New(trManager, receptionRepo, pvzLocker, productRepo)
	if err != nil {
		panic(fmt.Errorf("create product_adding use case: %w", err))
	}

	productRemoving, err := product_removing.New(trManager, receptionRepo, pvzLocker, productRepo)
	if err != nil {
		panic(fmt.Errorf("create product_removing use case: %w", err))
	}

	pvzListGetting, err := pvz_list_getting.New(pvzRepo, receptionRepo, productRepo)
	if err != nil {
		panic(fmt.Errorf("create pvz_list_getting use case: %w", err))
	}

	pvzRegistering, err := pvz_registering.New(pvzRepo)
	if err != nil {
		panic(fmt.Errorf("create pvz_registering use case: %w", err))
	}

	receptionClosing, err := reception_closing.New(trManager, receptionRepo, pvzLocker)
	if err != nil {
		panic(fmt.Errorf("create reception_closing use case: %w", err))
	}

	receptionCreating, err := reception_creating.New(trManager, receptionRepo, pvzLocker)
	if err != nil {
		panic(fmt.Errorf("create reception_creating use case: %w", err))
	}

	// Handlers
	dummyLoginHandler, err := dummy_login.New(dummyAuthentication, logger)
	if err != nil {
		panic(fmt.Errorf("create dummy_login handler: %w", err))
	}

	productAddHandler, err := product_add.New(productAdding, logger)
	if err != nil {
		panic(fmt.Errorf("create product_add handler: %w", err))
	}

	productRemoveLastHandler, err := product_remove_last.New(productRemoving, logger)
	if err != nil {
		panic(fmt.Errorf("create product_remove_last handler: %w", err))
	}

	pvzGetHandler, err := pvz_get.New(pvzListGetting, logger)
	if err != nil {
		panic(fmt.Errorf("create pvz_get handler: %w", err))
	}

	pvzRegisterHandler, err := pvz_register.New(pvzRegistering, logger)
	if err != nil {
		panic(fmt.Errorf("create pvz_register handler: %w", err))
	}

	receptionCloseHandler, err := reception_close.New(receptionClosing, logger)
	if err != nil {
		panic(fmt.Errorf("create reception_close handler: %w", err))
	}

	receptionCreateHandler, err := reception_create.New(receptionCreating, logger)
	if err != nil {
		panic(fmt.Errorf("create reception_create handler: %w", err))
	}

	noAuthMW, err := middleware.CreateNoAuthMiddleware()
	if err != nil {
		panic(fmt.Errorf("create no auth middleware: %w", err))
	}
	authMW, err := middleware.CreateAuthMiddleware(tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create auth middleware: %w", err))
	}

	authMux := http.NewServeMux()

	authMux.HandleFunc("POST /pvz", pvzRegisterHandler.Handle)
	authMux.HandleFunc("GET /pvz", pvzGetHandler.Handle)
	authMux.HandleFunc("POST /pvz/{pvzId}/close_last_reception", receptionCloseHandler.Handle)
	authMux.HandleFunc("POST /pvz/{pvzId}/delete_last_product", productRemoveLastHandler.Handle)
	authMux.HandleFunc("POST /receptions", receptionCreateHandler.Handle)
	authMux.HandleFunc("POST /products", productAddHandler.Handle)

	m := http.NewServeMux()
	m.Handle("POST /dummyLogin", noAuthMW(http.HandlerFunc(dummyLoginHandler.Handle)))
	m.Handle("/", authMW(authMux))

	s := &http.Server{
		Handler:           m,
		Addr:              "0.0.0.0:" + strconv.Itoa(cfg.ServerPort),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.Info("starting http server...")

	// And we serve HTTP until the world ends.
	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, context.Canceled) {
		panic(fmt.Errorf("http server ListenAndServe: %w", err))
	}
}
