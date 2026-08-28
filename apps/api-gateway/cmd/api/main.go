package main

import (
	"context"
	"os"

	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/auth"
	api "github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/handlers/router"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/api/services"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/config"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database/repository"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/database/seed"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/logger"
	"github.com/Samitdev0/cloudguard-sandbox/apps/api-gateway/internal/infrastructure/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()

	jwtService, err := auth.NewJWT(cfg.JWT)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to initialize JWT")
	}

	db, err := database.New(
		cfg.Database.DSN(),
		cfg.Database.MaxOpen,
		cfg.Database.MinOpen,
	)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().
			Err(err).
			Msg("database is unreachable")
	}

	log.Info().Msg("database connected")

	tenantRepository := repository.NewTenantRepository(db.Pool)
	userRepository := repository.NewUserRepository(db.Pool)
	auditRepository := repository.NewAuditRepository(db.Pool)

	if cfg.App.Environment == "development" {
		ctx := context.Background()

		if err := seed.Run(
			ctx,
			tenantRepository,
			userRepository,
		); err != nil {
			log.Fatal().
				Err(err).
				Msg("database seed failed")
		}

		log.Info().Msg("database seed completed")
	}

	authService := services.NewAuthService(
		userRepository,
		jwtService,
	)

	auditService := services.NewAuditService(
		auditRepository,
	)

	_ = auditService

	router := api.NewRouter(
		authService,
		jwtService,
		auditService,
	)

	srv := server.New(
		cfg.HTTP.Host,
		cfg.HTTP.Port,
		router,
		log,
	)

	if err := srv.Start(); err != nil {
		log.Error().
			Err(err).
			Msg("server stopped")

		os.Exit(1)
	}
}
