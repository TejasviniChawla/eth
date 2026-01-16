'use client';

export default function GlobalError({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-4 bg-canvas px-6 text-center">
      <h1 className="text-3xl font-semibold text-white">Something went wrong</h1>
      <p className="max-w-lg text-sm text-slate-300">{error.message}</p>
      <button
        type="button"
        onClick={() => reset()}
        className="rounded-full bg-accent px-5 py-2 text-sm font-semibold text-white"
      >
        Try again
      </button>
    </main>
  );
}
