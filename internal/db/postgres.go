package db

import (
	"context"
	"errors"
	"time"

	"github.com/etherfi/eth-dashboard/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) EnsureTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS wallet_snapshots (
			id SERIAL PRIMARY KEY,
			wallet_address TEXT NOT NULL,
			protocol TEXT NOT NULL,
			staked_amount NUMERIC NOT NULL,
			rewards NUMERIC NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS protocol_stats (
			id SERIAL PRIMARY KEY,
			protocol_name TEXT NOT NULL,
			current_apy NUMERIC NOT NULL,
			tvl NUMERIC NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS protocol_stats_protocol_name_idx ON protocol_stats (protocol_name);`,
	}

	for _, query := range queries {
		if _, err := p.pool.Exec(ctx, query); err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) UpsertProtocolStats(ctx context.Context, stats models.ProtocolStats) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO protocol_stats (protocol_name, current_apy, tvl, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (protocol_name) DO UPDATE
		SET current_apy = EXCLUDED.current_apy,
			tvl = EXCLUDED.tvl,
			updated_at = EXCLUDED.updated_at;
	`, stats.Protocol, stats.CurrentAPY, stats.TVL, stats.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetLatestProtocolStats(ctx context.Context) ([]models.ProtocolStats, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT protocol_name, current_apy, tvl, updated_at
		FROM protocol_stats
		ORDER BY protocol_name ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProtocolStats
	for rows.Next() {
		var stat models.ProtocolStats
		if err := rows.Scan(&stat.Protocol, &stat.CurrentAPY, &stat.TVL, &stat.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, stat)
	}
	return results, rows.Err()
}

func (p *Postgres) GetWalletSnapshots(ctx context.Context, address string) ([]models.YieldPoint, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT protocol, rewards, timestamp
		FROM wallet_snapshots
		WHERE wallet_address = $1
		ORDER BY timestamp ASC;
	`, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var yields []models.YieldPoint
	for rows.Next() {
		var point models.YieldPoint
		if err := rows.Scan(&point.Protocol, &point.Rewards, &point.Timestamp); err != nil {
			return nil, err
		}
		yields = append(yields, point)
	}
	return yields, rows.Err()
}

func (p *Postgres) InsertWalletSnapshot(ctx context.Context, address string, position models.StakingPosition) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO wallet_snapshots (wallet_address, protocol, staked_amount, rewards, timestamp)
		VALUES ($1, $2, $3, $4, $5);
	`, address, position.Protocol, position.StakedAmount, position.Rewards, time.Now().UTC())
	return err
}
