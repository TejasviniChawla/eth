package services

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/etherfi/eth-dashboard/internal/db"
	"github.com/etherfi/eth-dashboard/internal/models"
	"github.com/rs/zerolog/log"
)

type StakingService struct {
	store    *db.Postgres
	subgraph *SubgraphClient
}

func NewStakingService(store *db.Postgres, subgraph *SubgraphClient) *StakingService {
	return &StakingService{store: store, subgraph: subgraph}
}

func (s *StakingService) GetStakingPositions(ctx context.Context, address string) ([]models.StakingPosition, error) {
	normalized := strings.ToLower(address)

	etherfiResp, err := s.subgraph.FetchEtherFiStake(ctx, normalized)
	if err != nil {
		return nil, err
	}

	lidoResp, err := s.subgraph.FetchLidoStake(ctx, normalized)
	if err != nil {
		return nil, err
	}

	positions := []models.StakingPosition{
		mapEtherFiPosition(etherfiResp),
		mapLidoPosition(lidoResp),
		{
			Protocol:     models.ProtocolRocketPool,
			StakedAmount: 0,
			Rewards:      0,
			CurrentValue: 0,
		},
	}

	for _, position := range positions {
		if position.StakedAmount > 0 || position.Rewards > 0 {
			if err := s.store.InsertWalletSnapshot(ctx, normalized, position); err != nil {
				log.Error().Err(err).Str("protocol", string(position.Protocol)).Msg("failed to insert wallet snapshot")
			}
		}
	}

	return positions, nil
}

func (s *StakingService) GetYieldHistory(ctx context.Context, address string) ([]models.YieldPoint, error) {
	return s.store.GetWalletSnapshots(ctx, strings.ToLower(address))
}

func (s *StakingService) GetProtocolStats(ctx context.Context) ([]models.ProtocolStats, error) {
	return s.store.GetLatestProtocolStats(ctx)
}

func (s *StakingService) RefreshProtocolStats(ctx context.Context) {
	start := time.Now().Add(-24 * time.Hour).Unix()

	etherfiStats, err := s.buildProtocolStats(ctx, models.ProtocolEtherFi, s.subgraph.FetchEtherFiTVL, int(start))
	if err != nil {
		log.Error().Err(err).Msg("failed to refresh ether.fi stats")
	} else if err := s.store.UpsertProtocolStats(ctx, etherfiStats); err != nil {
		log.Error().Err(err).Msg("failed to upsert ether.fi stats")
	}

	lidoStats, err := s.buildProtocolStats(ctx, models.ProtocolLido, s.subgraph.FetchLidoTVL, int(start))
	if err != nil {
		log.Error().Err(err).Msg("failed to refresh lido stats")
	} else if err := s.store.UpsertProtocolStats(ctx, lidoStats); err != nil {
		log.Error().Err(err).Msg("failed to upsert lido stats")
	}

	rocketpoolStats := models.ProtocolStats{
		Protocol:  models.ProtocolRocketPool,
		CurrentAPY: 0,
		TVL:       0,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.store.UpsertProtocolStats(ctx, rocketpoolStats); err != nil {
		log.Error().Err(err).Msg("failed to upsert rocket pool stats")
	}
}

func (s *StakingService) buildProtocolStats(ctx context.Context, protocol models.Protocol, fetch func(context.Context, int) (protocolTVLResponse, error), startTime int) (models.ProtocolStats, error) {
	resp, err := fetch(ctx, startTime)
	if err != nil {
		return models.ProtocolStats{}, err
	}
	if len(resp.Data.ProtocolMetrics) == 0 {
		return models.ProtocolStats{Protocol: protocol, UpdatedAt: time.Now().UTC()}, nil
	}
	first := resp.Data.ProtocolMetrics[0]
	last := resp.Data.ProtocolMetrics[len(resp.Data.ProtocolMetrics)-1]

	startTVL, _ := strconv.ParseFloat(first.TotalValueLocked, 64)
	endTVL, _ := strconv.ParseFloat(last.TotalValueLocked, 64)

	apy := 0.0
	if startTVL > 0 {
		apy = ((endTVL - startTVL) / startTVL) * 365 * 100
	}

	return models.ProtocolStats{
		Protocol:  protocol,
		CurrentAPY: apy,
		TVL:       endTVL,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func mapEtherFiPosition(resp etherfiStakeResponse) models.StakingPosition {
	if len(resp.Data.Stakers) == 0 {
		return models.StakingPosition{Protocol: models.ProtocolEtherFi}
	}
	staker := resp.Data.Stakers[0]
	totalStaked, _ := strconv.ParseFloat(staker.TotalStaked, 64)
	totalWithdrawn, _ := strconv.ParseFloat(staker.TotalWithdrawn, 64)
	rewards := totalStaked - totalWithdrawn
	if rewards < 0 {
		rewards = 0
	}

	return models.StakingPosition{
		Protocol:     models.ProtocolEtherFi,
		StakedAmount: totalStaked,
		Rewards:      rewards,
		CurrentValue: totalStaked + rewards,
	}
}

func mapLidoPosition(resp lidoStakeResponse) models.StakingPosition {
	if len(resp.Data.Users) == 0 {
		return models.StakingPosition{Protocol: models.ProtocolLido}
	}
	user := resp.Data.Users[0]
	totalStaked, _ := strconv.ParseFloat(user.TotalStaked, 64)
	totalClaimed, _ := strconv.ParseFloat(user.TotalClaimed, 64)

	return models.StakingPosition{
		Protocol:     models.ProtocolLido,
		StakedAmount: totalStaked,
		Rewards:      totalClaimed,
		CurrentValue: totalStaked + totalClaimed,
	}
}

func StartProtocolStatsJob(ctx context.Context, service *StakingService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	service.RefreshProtocolStats(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			service.RefreshProtocolStats(ctx)
		}
	}
}
