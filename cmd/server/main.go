package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etherfi/eth-dashboard/internal/api"
	"github.com/etherfi/eth-dashboard/internal/db"
	"github.com/etherfi/eth-dashboard/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := envOrDefault("PORT", "8080")
	databaseURL := envOrDefault("DATABASE_URL", "")
	etherfiURL := envOrDefault("ETHERFI_SUBGRAPH_URL", "https://api.thegraph.com/subgraphs/name/etherfi-protocol/etherfi-mainnet")
	lidoURL := envOrDefault("LIDO_SUBGRAPH_URL", "https://api.thegraph.com/subgraphs/name/lidofinance/lido")

	pg, err := db.NewPostgres(ctx, databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer pg.Close()

	if err := pg.EnsureTables(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure tables")
	}

	subgraph := services.NewSubgraphClient(etherfiURL, lidoURL)
	stakingService := services.NewStakingService(pg, subgraph)

	go services.StartProtocolStatsJob(ctx, stakingService, 15*time.Minute)

	router := api.NewRouter(stakingService)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info().Msgf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
