/** @jsxImportSource solid-js */
import { secondaryButtonClass } from "../../../components/AuthCard";

export function AppBar(props: {
  organisation: string;
  project: string;
  onSignOut: () => void;
  onHelp?: () => void;
}) {
  return (
    <header class="flex items-center justify-between border-b border-paper-200 bg-paper px-4 py-3">
      <div class="flex items-center gap-3 text-sm font-medium text-ink">
        <span data-testid="org-name">{props.organisation}</span>
        <span class="text-ink-400">·</span>
        <span data-testid="project-name">{props.project}</span>
      </div>
      <div class="flex items-center gap-2">
        <button
          class={secondaryButtonClass}
          type="button"
          aria-label="Keyboard shortcuts"
          onClick={props.onHelp}
        >
          ?
        </button>
        <button
          class="rounded-full border border-paper-300 px-4 py-2 text-sm font-medium text-ink hover:bg-paper-50"
          type="button"
          onClick={props.onSignOut}
        >
          Sign out
        </button>
      </div>
    </header>
  );
}
