'use client';

import { useAccount } from 'wagmi';
import { useEffect, useMemo, useState } from 'react';
import WalletConnect from '../components/WalletConnect';
import StakingPositions from '../components/StakingPositions';
import YieldChart from '../components/YieldChart';
import ProtocolComparison from '../components/ProtocolComparison';
import type { ProtocolStats, StakingPosition, YieldPoint } from '@/types';
import { fetchProtocolStats, fetchStakingPositions, fetchYieldHistory } from '@/lib/api';

interface LoadState<T> {
  data: T;
  isLoading: boolean;
  error: string | null;
}

const initialState = <T,>(value: T): LoadState<T> => ({
  data: value,
  isLoading: true,
  error: null
});

export default function DashboardPage() {
  const { address, isConnected } = useAccount();

  const [positionsState, setPositionsState] = useState<LoadState<StakingPosition[]>>(
    initialState<StakingPosition[]>([])
  );
  const [yieldState, setYieldState] = useState<LoadState<YieldPoint[]>>(initialState<YieldPoint[]>([]));
  const [protocolState, setProtocolState] = useState<LoadState<ProtocolStats[]>>(
    initialState<ProtocolStats[]>([])
  );

  useEffect(() => {
    if (!address) {
      return;
    }

    let mounted = true;

    const loadData = async () => {
      try {
        const [positions, yieldHistory, protocolStats] = await Promise.all([
          fetchStakingPositions(address),
          fetchYieldHistory(address),
          fetchProtocolStats()
        ]);

        if (mounted) {
          setPositionsState({ data: positions, isLoading: false, error: null });
          setYieldState({ data: yieldHistory, isLoading: false, error: null });
          setProtocolState({ data: protocolStats, isLoading: false, error: null });
        }
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to load data.';
        if (mounted) {
          setPositionsState((prev) => ({ ...prev, isLoading: false, error: message }));
          setYieldState((prev) => ({ ...prev, isLoading: false, error: message }));
          setProtocolState((prev) => ({ ...prev, isLoading: false, error: message }));
        }
      }
    };

    loadData();

    return () => {
      mounted = false;
    };
  }, [address]);

  const totalRewards = useMemo(
    () => positionsState.data.reduce((sum, position) => sum + position.rewards, 0),
    [positionsState.data]
  );

  if (!isConnected) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center gap-6 bg-canvas px-6 py-12">
        <h1 className="text-3xl font-semibold text-white">Connect your wallet to view the dashboard.</h1>
        <WalletConnect />
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-canvas px-6 py-10">
      <div className="mx-auto flex max-w-6xl flex-col gap-8">
        <header className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-slate-400">Dashboard</p>
            <h1 className="text-3xl font-semibold text-white">Staking overview</h1>
          </div>
          <WalletConnect />
        </header>

        <section className="grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-card p-5">
            <p className="text-sm text-slate-400">Connected Wallet</p>
            <p className="mt-2 text-lg text-white">
              {address?.slice(0, 6)}...{address?.slice(-4)}
            </p>
          </div>
          <div className="rounded-2xl border border-white/10 bg-card p-5">
            <p className="text-sm text-slate-400">Total Rewards</p>
            <p className="mt-2 text-lg text-white">{totalRewards.toFixed(4)} ETH</p>
          </div>
          <div className="rounded-2xl border border-white/10 bg-card p-5">
            <p className="text-sm text-slate-400">Last Updated</p>
            <p className="mt-2 text-lg text-white">Just now</p>
          </div>
        </section>

        <section>
          <StakingPositions
            positions={positionsState.data}
            isLoading={positionsState.isLoading}
            error={positionsState.error}
          />
        </section>

        <section className="grid gap-6 lg:grid-cols-2">
          <YieldChart data={yieldState.data} isLoading={yieldState.isLoading} error={yieldState.error} />
          <ProtocolComparison
            stats={protocolState.data}
            positions={positionsState.data}
            isLoading={protocolState.isLoading}
            error={protocolState.error}
          />
        </section>
      </div>
    </main>
  );
}
