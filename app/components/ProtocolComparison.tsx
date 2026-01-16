import type { ProtocolStats, StakingPosition } from '@/types';

interface ProtocolComparisonProps {
  stats: ProtocolStats[];
  positions: StakingPosition[];
  isLoading: boolean;
  error?: string | null;
}

export default function ProtocolComparison({ stats, positions, isLoading, error }: ProtocolComparisonProps) {
  if (isLoading) {
    return <div className="h-48 animate-pulse rounded-2xl border border-white/10 bg-card" />;
  }

  if (error) {
    return <div className="rounded-xl border border-red-500/40 bg-red-500/10 p-4 text-sm">{error}</div>;
  }

  const positionMap = new Map(positions.map((position) => [position.protocol, position]));

  return (
    <div className="rounded-2xl border border-white/10 bg-card p-5">
      <h3 className="mb-4 text-lg font-semibold text-white">Protocol Comparison</h3>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm">
          <thead className="text-xs uppercase tracking-wide text-slate-400">
            <tr>
              <th className="pb-3">Protocol</th>
              <th className="pb-3">APY</th>
              <th className="pb-3">TVL</th>
              <th className="pb-3">Your Stake</th>
            </tr>
          </thead>
          <tbody className="text-slate-200">
            {stats.map((stat) => {
              const position = positionMap.get(stat.protocol);
              return (
                <tr key={stat.protocol} className="border-t border-white/5">
                  <td className="py-3 capitalize">{stat.protocol}</td>
                  <td className="py-3">{stat.currentApy.toFixed(2)}%</td>
                  <td className="py-3">{stat.tvl.toLocaleString()} ETH</td>
                  <td className="py-3">{position ? `${position.stakedAmount.toFixed(2)} ETH` : '—'}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
