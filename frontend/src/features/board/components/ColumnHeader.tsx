/** @jsxImportSource solid-js */
import { cn } from "../../../lib/cn";

export function ColumnHeader(props: { title: string; current?: boolean; extra?: string }) {
  return (
    <header
      class={cn(
        "border-b border-paper-300 px-4 py-3 text-xs font-medium uppercase tracking-wider text-ink-700",
        props.current && "bg-bloom text-white",
      )}
    >
      <div>{props.title}</div>
      {props.extra && <div class="mt-1 text-xs font-medium normal-case tracking-normal">{props.extra}</div>}
    </header>
  );
}
