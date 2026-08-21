/** @jsxImportSource solid-js */
import type { JSX } from "solid-js";

export function AuthCard(props: { title: string; children: JSX.Element }) {
  return (
    <main class="mx-auto flex min-h-screen max-w-md flex-col justify-center gap-6 px-6 py-16">
      <p class="text-xs font-medium uppercase tracking-[0.2em] text-bloom">Flower</p>
      <h1 class="font-heading text-3xl font-semibold text-ink">{props.title}</h1>
      {props.children}
    </main>
  );
}

export const inputClass =
  "w-full rounded-lg border border-paper-300 px-3 py-2 text-sm focus:border-bloom focus:outline-none focus:ring-1 focus:ring-bloom";

export const labelClass = "text-sm font-medium text-ink-800";

export const primaryButtonClass =
  "rounded-full bg-bloom px-4 py-2 text-sm font-medium text-white hover:bg-bloom/90 disabled:cursor-not-allowed disabled:opacity-50";

export const secondaryButtonClass =
  "rounded-full border border-paper-300 px-4 py-2 text-sm font-medium text-ink hover:bg-paper-50 disabled:cursor-not-allowed disabled:opacity-50";
