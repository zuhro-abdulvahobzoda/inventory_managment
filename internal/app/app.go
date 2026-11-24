package app

import (
	"context"
	"errors"
	"fmt"
	"inventory-api/internal/config"
	"inventory-api/internal/handler"
	itemrepo "inventory-api/internal/repository/item"
	locationrepo "inventory-api/internal/repository/location"
	userrepo "inventory-api/internal/repository/user"
	"inventory-api/internal/service"
	"inventory-api/internal/service/auth"
	"inventory-api/internal/service/item"
	"inventory-api/internal/service/location"
	jwttoken "inventory-api/internal/service/token/jwt"
	"inventory-api/internal/service/user"
	"inventory-api/migrations"
	"inventory-api/pkg/db/postgres"
	"inventory-api/pkg/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slog.SetDefault(logger.New(cfg.Logger))

	if cfg.ServerConfigs.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := postgres.Migrate(cfg.Postgres, migrations.FS, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	pool, err := postgres.New(context.Background(), cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()

	userRepo := userrepo.NewUserRepo(pool)
	itemRepo := itemrepo.NewItemRepository(pool)
	locationRepo := locationrepo.NewLocationRepository(pool)

	tokens := jwttoken.NewTokenManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	srvs := &service.Services{
		AuthService:     auth.NewAuthService(userRepo, tokens),
		UserService:     user.NewUserService(userRepo),
		ItemService:     item.NewItemService(itemRepo),
		LocationService: location.NewLocationService(locationRepo),
	}
	h := handler.NewHandler(srvs)

	server := &http.Server{
		Addr:              ":" + cfg.ServerConfigs.Port,
		Handler:           h.InitRoutes(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Lauching server in goroutine so main thread could keep "listening" to stop signal
	go func() {
		log.Printf("http server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	// Waits Ctrl+C (SIGINT) or SIGTERM from Docker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	log.Println("server stopped")
	return nil
}
