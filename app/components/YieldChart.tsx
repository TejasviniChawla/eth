'use client';

import type { YieldPoint } from '@/types';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid
} from 'recharts';

interface YieldChartProps {
  data: YieldPoint[];
  isLoading: boolean;
  error?: string | null;
}

export default function YieldChart({ data, isLoading, error }: YieldChartProps) {
  if (isLoading) {
    return <div className="h-64 animate-pulse rounded-2xl border border-white/10 bg-card" />;
  }

  if (error) {
    return <div className="rounded-xl border border-red-500/40 bg-red-500/10 p-4 text-sm">{error}</div>;
  }

  return (
    <div className="h-80 rounded-2xl border border-white/10 bg-card p-5">
      <h3 className="mb-4 text-lg font-semibold text-white">Yield Over Time</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ top: 10, right: 20, left: 0, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#1F2937" />
          <XAxis
            dataKey="timestamp"
            tick={{ fill: '#94A3B8', fontSize: 12 }}
            tickFormatter={(value) => new Date(value as string).toLocaleDateString()}
          />
          <YAxis tick={{ fill: '#94A3B8', fontSize: 12 }} />
          <Tooltip
            contentStyle={{
              background: '#0F172A',
              border: '1px solid rgba(255,255,255,0.1)'
            }}
          />
          <Line type="monotone" dataKey="rewards" stroke="#8B5CF6" strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
