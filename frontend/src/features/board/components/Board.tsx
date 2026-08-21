/** @jsxImportSource solid-js */
import { Show, createSignal, onCleanup, onMount } from "solid-js";
import { AppBar } from "./AppBar";
import { BoardColumn } from "./BoardColumn";
import { primaryButtonClass, secondaryButtonClass } from "../../../components/AuthCard";
import type { Pack } from "../../../lib/api/stories";

function formatEnds(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleDateString("en-GB", { day: "numeric", month: "short", timeZone: "Australia/Melbourne" });
}

function isTypingTarget(event: KeyboardEvent): boolean {
  const target = event.target as HTMLElement | null;
  return Boolean(target && target.closest("input, textarea, select, [contenteditable=true]"));
}

export function Board(props: {
  organisation: string;
  project: string;
  pack: Pack;
  onSignOut: () => void;
}) {
  const extra = `${props.pack.current_points} / ${props.pack.denominator} · Ends ${formatEnds(props.pack.current_window_ends_at)}`;
  const [help, setHelp] = createSignal(false);
  let addButton: HTMLButtonElement | undefined;

  function onKey(event: KeyboardEvent) {
    if (isTypingTarget(event)) {
      return;
    }
    if (event.key === "Escape") {
      setHelp(false);
      return;
    }
    if (event.key === "?") {
      event.preventDefault();
      setHelp(true);
      return;
    }
    if ((event.key === "a" || event.key === "A") && !event.metaKey && !event.ctrlKey && !event.altKey) {
      event.preventDefault();
      addButton?.click();
    }
  }

  onMount(() => {
    document.addEventListener("keydown", onKey);
    onCleanup(() => document.removeEventListener("keydown", onKey));
  });

  return (
    <div class="flex h-screen flex-col bg-paper" data-testid="board" tabIndex={-1} onKeyDown={onKey}>
      <AppBar organisation={props.organisation} project={props.project} onSignOut={props.onSignOut} onHelp={() => setHelp(true)} />
      <div class="flex min-h-0 flex-1 overflow-x-auto">
        <BoardColumn
          title="Icebox"
          empty="Nothing waiting. Capture a story."
          action={
            <button
              ref={addButton}
              class={primaryButtonClass}
              type="button"
              data-testid="add-story"
            >
              Add a story
            </button>
          }
        />
        <BoardColumn title="Backlog" empty="Nothing ranked yet. Pull from Icebox when you’re ready." />
        <BoardColumn
          title="Current"
          current
          extra={extra}
          empty="Nothing in this window yet. Pull a story from Icebox or create one."
        />
        <BoardColumn
          title="Done"
          empty="Nothing in Done yet. Accepted work lands here after it ages out of this window."
        />
      </div>
      <Show when={help()}>
        <div class="fixed inset-0 z-50 flex items-center justify-center bg-ink/40" role="dialog" aria-label="Keyboard shortcuts" aria-modal="true">
          <div class="rounded-lg border border-paper-300 bg-white p-6">
            <p class="font-heading text-xl font-semibold text-ink">Shortcuts</p>
            <button class={`${secondaryButtonClass} mt-4`} type="button" onClick={() => setHelp(false)}>
              Close
            </button>
          </div>
        </div>
      </Show>
    </div>
  );
}
