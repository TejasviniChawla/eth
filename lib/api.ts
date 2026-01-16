import type { ProtocolStats, StakingPosition, YieldPoint } from '@/types';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? 'http://localhost:8080';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || 'Request failed');
  }
  return (await response.json()) as T;
}

export async function fetchStakingPositions(address: string): Promise<StakingPosition[]> {
  const response = await fetch(`${API_BASE}/api/staking/${address}`, {
    cache: 'no-store'
  });
  return handleResponse<StakingPosition[]>(response);
}

export async function fetchYieldHistory(address: string): Promise<YieldPoint[]> {
  const response = await fetch(`${API_BASE}/api/yields/${address}`, {
    cache: 'no-store'
  });
  return handleResponse<YieldPoint[]>(response);
}

export async function fetchProtocolStats(): Promise<ProtocolStats[]> {
  const response = await fetch(`${API_BASE}/api/protocols`, {
    cache: 'no-store'
  });
  return handleResponse<ProtocolStats[]>(response);
}
