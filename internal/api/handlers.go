package api

import (
	"encoding/json"
	"net/http"

	"github.com/etherfi/eth-dashboard/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	stakingService *services.StakingService
}

func NewHandler(stakingService *services.StakingService) *Handler {
	return &Handler{stakingService: stakingService}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetStaking(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	positions, err := h.stakingService.GetStakingPositions(r.Context(), address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("failed to get staking positions")
		respondError(w, http.StatusBadGateway, "failed to fetch staking positions")
		return
	}
	respondJSON(w, http.StatusOK, positions)
}

func (h *Handler) GetYields(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	yields, err := h.stakingService.GetYieldHistory(r.Context(), address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("failed to get yield history")
		respondError(w, http.StatusBadGateway, "failed to fetch yield history")
		return
	}
	respondJSON(w, http.StatusOK, yields)
}

func (h *Handler) GetProtocols(w http.ResponseWriter, r *http.Request) {
	protocols, err := h.stakingService.GetProtocolStats(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get protocol stats")
		respondError(w, http.StatusBadGateway, "failed to fetch protocol stats")
		return
	}
	respondJSON(w, http.StatusOK, protocols)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error().Err(err).Msg("failed to encode json response")
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
