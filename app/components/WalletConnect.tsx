'use client';

import { useAccount, useConnect, useDisconnect, useEnsName } from 'wagmi';
import { useMemo } from 'react';

export default function WalletConnect() {
  const { address, isConnected } = useAccount();
  const { data: ensName } = useEnsName({ address, chainId: 1 });
  const { connectors, connect, isPending } = useConnect();
  const { disconnect } = useDisconnect();

  const displayAddress = useMemo(() => {
    if (!address) {
      return '';
    }
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  }, [address]);

  if (isConnected && address) {
    return (
      <div className="flex items-center gap-3 rounded-full border border-white/10 bg-card px-4 py-2">
        <span className="text-sm text-slate-200">
          {ensName ? `${ensName} (${displayAddress})` : displayAddress}
        </span>
        <button
          type="button"
          onClick={() => disconnect()}
          className="rounded-full bg-white/10 px-3 py-1 text-xs uppercase tracking-wide text-slate-200"
        >
          Disconnect
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap gap-3">
      {connectors.map((connector) => (
        <button
          key={connector.uid}
          type="button"
          disabled={!connector.ready || isPending}
          onClick={() => connect({ connector })}
          className="rounded-full bg-accent px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-purple-500/30 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isPending ? 'Connecting...' : `Connect ${connector.name}`}
        </button>
      ))}
    </div>
  );
}
