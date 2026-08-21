/** @jsxImportSource solid-js */
import type { JSX } from "solid-js";
import { ColumnHeader } from "./ColumnHeader";

export function BoardColumn(props: {
  title: string;
  current?: boolean;
  extra?: string;
  empty: string;
  action?: JSX.Element;
}) {
  return (
    <section class="flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper" data-testid={`column-${props.title.toLowerCase()}`}>
      <ColumnHeader title={props.title} current={props.current} extra={props.extra} />
      <div class="flex flex-1 flex-col gap-4 p-4">
        <p class="text-sm text-ink-700">{props.empty}</p>
        {props.action}
      </div>
    </section>
  );
}
