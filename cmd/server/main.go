package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-pvz/internal/api/dummy_login"
	"github.com/inna-maikut/avito-pvz/internal/api/login"
	"github.com/inna-maikut/avito-pvz/internal/api/product_add"
	"github.com/inna-maikut/avito-pvz/internal/api/product_remove_last"
	"github.com/inna-maikut/avito-pvz/internal/api/pvz_get"
	"github.com/inna-maikut/avito-pvz/internal/api/pvz_register"
	"github.com/inna-maikut/avito-pvz/internal/api/reception_close"
	"github.com/inna-maikut/avito-pvz/internal/api/reception_create"
	"github.com/inna-maikut/avito-pvz/internal/api/register"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/config"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/metrics"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/middleware"
	"github.com/inna-maikut/avito-pvz/internal/infrastructure/pg"
	"github.com/inna-maikut/avito-pvz/internal/repository"
	"github.com/inna-maikut/avito-pvz/internal/usecases/authenticating"
	"github.com/inna-maikut/avito-pvz/internal/usecases/dummy_authenticating"
	"github.com/inna-maikut/avito-pvz/internal/usecases/product_adding"
	"github.com/inna-maikut/avito-pvz/internal/usecases/product_removing"
	"github.com/inna-maikut/avito-pvz/internal/usecases/pvz_list_getting"
	"github.com/inna-maikut/avito-pvz/internal/usecases/pvz_registering"
	"github.com/inna-maikut/avito-pvz/internal/usecases/reception_closing"
	"github.com/inna-maikut/avito-pvz/internal/usecases/reception_creating"
	"github.com/inna-maikut/avito-pvz/internal/usecases/registering"
)

const (
	readHeaderTimeout = time.Second
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer func() {
		_ = logger.Sync()
	}()

	// Graceful shutdown - cancel context on signals SIGINT and SIGTERM
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("stopping...")

		cancel()
	}()

	// catch init panics
	defer func() {
		if panicErr := recover(); panicErr != nil {
			if typedErr, ok := panicErr.(error); ok {
				logger.Error("panic error", zap.Error(typedErr))
			} else {
				logger.Error("panic", zap.Any("message", panicErr))
			}
		}
	}()

	metric, err := metrics.New()
	if err != nil {
		panic(fmt.Errorf("create metrics: %w", err))
	}

	// Postgres DB and repositories

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

	userRepo, err := repository.NewUserRepository(db, trmsqlx.DefaultCtxGetter)
	if err != nil {
		panic(fmt.Errorf("create user repository: %w", err))
	}

	// Use cases

	dummyAuthentication, err := dummy_authenticating.New(tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create dummy_authenticating use case: %w", err))
	}

	authentication, err := authenticating.New(userRepo, tokenProvider)
	if err != nil {
		panic(fmt.Errorf("create authenticating use case: %w", err))
	}

	registration, err := registering.New(userRepo)
	if err != nil {
		panic(fmt.Errorf("create registering use case: %w", err))
	}

	productAdding, err := product_adding.New(trManager, receptionRepo, pvzLocker, productRepo, metric)
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

	pvzRegistering, err := pvz_registering.New(pvzRepo, metric)
	if err != nil {
		panic(fmt.Errorf("create pvz_registering use case: %w", err))
	}

	receptionClosing, err := reception_closing.New(trManager, receptionRepo, pvzLocker)
	if err != nil {
		panic(fmt.Errorf("create reception_closing use case: %w", err))
	}

	receptionCreating, err := reception_creating.New(trManager, receptionRepo, pvzLocker, metric)
	if err != nil {
		panic(fmt.Errorf("create reception_creating use case: %w", err))
	}

	// API Handlers

	dummyLoginHandler, err := dummy_login.New(dummyAuthentication, logger)
	if err != nil {
		panic(fmt.Errorf("create dummy_login handler: %w", err))
	}

	loginHandler, err := login.New(authentication, logger)
	if err != nil {
		panic(fmt.Errorf("create login handler: %w", err))
	}

	registerHandler, err := register.New(registration, logger)
	if err != nil {
		panic(fmt.Errorf("create register handler: %w", err))
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

	// HTTP server set up

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
	m.Handle("POST /register", noAuthMW(http.HandlerFunc(loginHandler.Handle)))
	m.Handle("POST /login", noAuthMW(http.HandlerFunc(registerHandler.Handle)))
	m.Handle("/", authMW(authMux))
	handler := metric.HTTPServerMW(m)

	var wg sync.WaitGroup

	// metrics http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		metric.RunHTTPServer(ctx, cfg, logger)
	}()

	// http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		runHTTPServer(ctx, handler, cfg, logger)
	}()

	wg.Wait()
	logger.Info("successful stop")
}

func runHTTPServer(ctx context.Context, handler http.Handler, cfg config.Config, logger *zap.Logger) {
	s := &http.Server{
		Handler:           handler,
		Addr:              cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if shutdownErr := s.Shutdown(shutdownCtx); shutdownErr != nil {
			shutdownErr = fmt.Errorf("HTTP shutdown error: %w", shutdownErr)
			logger.Error("HTTP shutdown error", zap.Error(shutdownErr))
		}
	}()

	logger.Info("starting http server...")

	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		err = fmt.Errorf("HTTP server ListenAndServe: %w", err)
		logger.Error("HTTP server ListenAndServe", zap.Error(err))
	}
}
