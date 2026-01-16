import Link from 'next/link';
import WalletConnect from './components/WalletConnect';

export default function HomePage() {
  return (
    <main className="min-h-screen bg-canvas px-6 py-12">
      <div className="mx-auto flex max-w-5xl flex-col gap-10">
        <header className="flex flex-col gap-6">
          <p className="text-sm uppercase tracking-[0.3em] text-slate-400">ether.fi staking dashboard</p>
          <h1 className="text-4xl font-semibold text-white md:text-6xl">
            One dashboard to track <span className="gradient-text">staking yield</span> across protocols.
          </h1>
          <p className="max-w-2xl text-lg text-slate-300">
            Connect your wallet to compare ether.fi, Lido, and Rocket Pool positions, monitor yield over
            time, and stay ahead of protocol APY movements.
          </p>
          <WalletConnect />
          <Link
            href="/dashboard"
            className="w-fit rounded-full border border-white/10 px-5 py-2 text-sm uppercase tracking-wide text-slate-200 hover:border-white/30"
          >
            View Dashboard
          </Link>
        </header>
        <section className="grid gap-6 md:grid-cols-3">
          {[
            {
              title: 'Unified Positions',
              description: 'Aggregated staking positions across the top LSD protocols.'
            },
            {
              title: 'Yield Tracking',
              description: 'Visualize reward accumulation over time with subgraph data.'
            },
            {
              title: 'Protocol Insights',
              description: 'Compare APY and TVL updates every 15 minutes.'
            }
          ].map((item) => (
            <div key={item.title} className="rounded-2xl border border-white/10 bg-card p-6">
              <h3 className="text-lg font-semibold text-white">{item.title}</h3>
              <p className="mt-2 text-sm text-slate-300">{item.description}</p>
            </div>
          ))}
        </section>
      </div>
    </main>
  );
}
