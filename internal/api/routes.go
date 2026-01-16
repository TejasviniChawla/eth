package api

import (
	"net/http"
	"time"

	"github.com/etherfi/eth-dashboard/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(stakingService *services.StakingService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	handler := NewHandler(stakingService)

	r.Get("/api/health", handler.Health)
	r.Get("/api/staking/{address}", handler.GetStaking)
	r.Get("/api/yields/{address}", handler.GetYields)
	r.Get("/api/protocols", handler.GetProtocols)

	return r
}
