import { http, createConfig, createStorage } from 'wagmi';
import { mainnet } from 'wagmi/chains';
import { injected, walletConnect } from 'wagmi/connectors';

const walletConnectProjectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID ?? '';

export const wagmiConfig = createConfig({
  chains: [mainnet],
  connectors: [
    injected({
      target: 'metaMask'
    }),
    walletConnect({
      projectId: walletConnectProjectId,
      showQrModal: true
    })
  ],
  transports: {
    [mainnet.id]: http()
  },
  ssr: true,
  storage: createStorage({
    storage: typeof window !== 'undefined' ? window.localStorage : undefined
  })
});
