/** @jsxImportSource solid-js */
import { AppBar } from "./AppBar";
import { BoardColumn } from "./BoardColumn";
import { primaryButtonClass } from "../../../components/AuthCard";
import type { Pack } from "../../../lib/api/stories";

function formatEnds(iso: string): string {
  const date = new Date(iso);
  return date.toLocaleDateString("en-GB", { day: "numeric", month: "short", timeZone: "Australia/Melbourne" });
}

export function EmptyBoard(props: {
  organisation: string;
  project: string;
  pack: Pack;
  onSignOut: () => void;
}) {
  const extra = `${props.pack.current_points} / ${props.pack.denominator} · Ends ${formatEnds(props.pack.current_window_ends_at)}`;
  return (
    <div class="flex h-screen flex-col bg-paper">
      <AppBar organisation={props.organisation} project={props.project} onSignOut={props.onSignOut} />
      <div class="flex min-h-0 flex-1 overflow-x-auto">
        <BoardColumn
          title="Icebox"
          empty="Nothing waiting. Capture a story."
          action={
            <button class={primaryButtonClass} type="button">
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
    </div>
  );
}
