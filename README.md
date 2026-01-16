# ether.fi Staking Dashboard & Yield Tracker

A production-ready, full-stack dashboard that connects to a user’s wallet, aggregates staking positions across ether.fi, Lido, and Rocket Pool, tracks yield history, and compares protocol APY with subgraph-powered data.

## ✨ Highlights

- **Wallet-aware dashboard** with wagmi + viem connectors (MetaMask + WalletConnect).
- **Real subgraph data** from The Graph hosted endpoints for ether.fi and Lido.
- **PostgreSQL-backed snapshots** for historical yield tracking.
- **15-minute background job** that caches protocol stats.
- **Dark-mode UI** built with Tailwind + Next.js App Router.

---

## Architecture

```
┌──────────────────────────────┐
│         Next.js UI           │
│  WalletConnect + Dashboard   │
└──────────────┬───────────────┘
               │ REST
┌──────────────▼───────────────┐
│        Go API (Chi)          │
│  /api/staking /api/yields    │
│  /api/protocols /api/health  │
└──────────────┬───────────────┘
               │ SQL + Subgraphs
┌──────────────▼───────────────┐
│     PostgreSQL + The Graph   │
└──────────────────────────────┘
```

---

## Backend (Go + Chi)

**Project structure**

```
/cmd/server/main.go
/internal/api/handlers.go
/internal/api/routes.go
/internal/services/staking.go
/internal/services/subgraph.go
/internal/models/types.go
/internal/db/postgres.go
/pkg/ethereum/client.go
```

### API Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/staking/:address` | Staking positions for a wallet |
| GET | `/api/yields/:address` | Historical yield data |
| GET | `/api/protocols` | Supported protocols with APY |

### Environment Variables

| Variable | Description |
|---------|-------------|
| `PORT` | API port (default: 8080) |
| `DATABASE_URL` | Postgres connection string |
| `ETHERFI_SUBGRAPH_URL` | Ether.fi subgraph endpoint |
| `LIDO_SUBGRAPH_URL` | Lido subgraph endpoint |

### Database Schema

```sql
CREATE TABLE wallet_snapshots (
  id SERIAL PRIMARY KEY,
  wallet_address TEXT NOT NULL,
  protocol TEXT NOT NULL,
  staked_amount NUMERIC NOT NULL,
  rewards NUMERIC NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL
);

CREATE TABLE protocol_stats (
  id SERIAL PRIMARY KEY,
  protocol_name TEXT NOT NULL UNIQUE,
  current_apy NUMERIC NOT NULL,
  tvl NUMERIC NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

### Background Job

A goroutine (ticker every 15 minutes) runs `RefreshProtocolStats` to update cached APY and TVL data.

---

## Subgraph Queries (Implemented)

```graphql
# Query for ether.fi staking positions
query GetUserStakes($user: String!) {
  stakers(where: { id: $user }) {
    id
    totalStaked
    totalWithdrawn
    stakes {
      amount
      timestamp
      validator {
        id
        status
      }
    }
  }
}

# Query for protocol TVL over time
query GetProtocolTVL($startTime: Int!) {
  protocolMetrics(
    where: { timestamp_gt: $startTime }
    orderBy: timestamp
    orderDirection: asc
  ) {
    timestamp
    totalValueLocked
    totalStakers
  }
}
```

---

## Frontend (Next.js 14 + TypeScript)

**Project structure**

```
/app/page.tsx
/app/dashboard/page.tsx
/app/components/WalletConnect.tsx
/app/components/StakingPositions.tsx
/app/components/YieldChart.tsx
/app/components/ProtocolComparison.tsx
/lib/api.ts
/lib/wagmi.ts
/types/index.ts
```

### Features

- Wallet connection with ENS resolution
- Dashboard gated by wallet connection
- Protocol comparison table (APY, TVL, personal position)
- Yield chart (Recharts)
- Tailwind styling with dark mode

### Environment Variables

| Variable | Description |
|---------|-------------|
| `NEXT_PUBLIC_API_BASE` | Backend base URL |
| `NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID` | WalletConnect project ID |

---

## Local Development

### 1) Backend

```bash
export DATABASE_URL=postgres://etherfi:etherfi@localhost:5432/etherfi?sslmode=disable
export ETHERFI_SUBGRAPH_URL=https://api.thegraph.com/subgraphs/name/etherfi-protocol/etherfi-mainnet
export LIDO_SUBGRAPH_URL=https://api.thegraph.com/subgraphs/name/lidofinance/lido

go run ./cmd/server
```

### 2) Frontend

```bash
export NEXT_PUBLIC_API_BASE=http://localhost:8080
export NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=your_project_id

npm install
npm run dev
```

---

## Docker & Compose

### Build images

```bash
docker build -f Dockerfile.backend -t etherfi-backend .
docker build -f Dockerfile.frontend -t etherfi-frontend .
```

### Compose (Backend + Frontend + Postgres)

```bash
docker compose up --build
```

---

## CI/CD

The GitHub Actions workflow (`.github/workflows/ci.yml`) performs:

- Go tests
- TypeScript type checking
- Docker image builds

---

## Production Notes

- The backend uses a static build (CGO disabled) for a minimal runtime image.
- Ensure `NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID` is set for WalletConnect to work.
- The Graph endpoints are public; add rate limiting for production usage.

---

## Roadmap

- Add Rocket Pool subgraph integration for real data.
- Persist wallet snapshots on a scheduled job.
- Improve APY calculation with protocol-specific metrics.
