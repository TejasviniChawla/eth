import type { StakingPosition } from '@/types';

interface StakingPositionsProps {
  positions: StakingPosition[];
  isLoading: boolean;
  error?: string | null;
}

export default function StakingPositions({ positions, isLoading, error }: StakingPositionsProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {Array.from({ length: 3 }).map((_, index) => (
          <div
            key={`skeleton-${index}`}
            className="h-32 animate-pulse rounded-2xl border border-white/10 bg-card"
          />
        ))}
      </div>
    );
  }

  if (error) {
    return <div className="rounded-xl border border-red-500/40 bg-red-500/10 p-4 text-sm">{error}</div>;
  }

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {positions.map((position) => (
        <div key={position.protocol} className="rounded-2xl border border-white/10 bg-card p-5">
          <p className="text-sm uppercase tracking-wide text-slate-400">{position.protocol}</p>
          <h3 className="mt-2 text-2xl font-semibold text-white">
            {position.stakedAmount.toFixed(2)} ETH
          </h3>
          <p className="mt-2 text-sm text-slate-300">
            Rewards: <span className="text-slate-100">{position.rewards.toFixed(4)} ETH</span>
          </p>
          <p className="text-sm text-slate-300">
            Current Value:{' '}
            <span className="text-slate-100">{position.currentValue.toFixed(4)} ETH</span>
          </p>
        </div>
      ))}
    </div>
  );
}
