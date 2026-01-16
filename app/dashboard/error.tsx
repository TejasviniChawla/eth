'use client';

export default function DashboardError({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div className="rounded-2xl border border-red-500/40 bg-red-500/10 p-6 text-sm text-slate-100">
      <p className="mb-3 font-semibold">Dashboard failed to load.</p>
      <p className="mb-4 text-slate-200">{error.message}</p>
      <button
        type="button"
        onClick={() => reset()}
        className="rounded-full bg-accent px-4 py-2 text-xs font-semibold uppercase tracking-wide text-white"
      >
        Retry
      </button>
    </div>
  );
}
