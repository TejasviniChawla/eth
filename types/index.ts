export type Protocol = 'ether.fi' | 'lido' | 'rocketpool';

export interface StakingPosition {
  protocol: Protocol;
  stakedAmount: number;
  rewards: number;
  currentValue: number;
}

export interface YieldPoint {
  timestamp: string;
  rewards: number;
  protocol: Protocol;
}

export interface ProtocolStats {
  protocol: Protocol;
  currentApy: number;
  tvl: number;
  updatedAt: string;
}
