/** @jsxImportSource solid-js */
import { For, Show, createResource } from "solid-js";
import { fetchHealth } from "../lib/api";
import { cn } from "../lib/cn";

const PANELS = [
  { title: "Icebox", body: "Unscheduled stories wait here until they are ranked." },
  { title: "Backlog", body: "A single ranked queue of work that will be scheduled next." },
  { title: "Current", body: "The iteration in flight, paced by accepted points." },
  { title: "Done", body: "Accepted stories from recent iterations." },
] as const;

type HealthView =
  | { kind: "ok"; version: string; status: string }
  | { kind: "error" };

async function loadHealth(): Promise<HealthView> {
  try {
    const health = await fetchHealth();
    return { kind: "ok", version: health.version, status: health.status };
  } catch (error) {
    if (!(error instanceof Error)) {
      throw error;
    }
    return { kind: "error" };
  }
}

function ApiStatusMessage(props: { view: HealthView }) {
  if (props.view.kind === "error") {
    return (
      <p class="mt-2 text-sm text-red-700" data-testid="api-status-error">
        The API is not reachable.
      </p>
    );
  }

  return (
    <p class="mt-2 text-sm text-stem" data-testid="api-status-ok">
      Connected to Flower API {props.view.version} ({props.view.status}).
    </p>
  );
}

export function HomePage() {
  const [health] = createResource(loadHealth);

  return (
    <main class="mx-auto flex min-h-screen max-w-3xl flex-col justify-center gap-8 px-6 py-16">
      <div class="flex flex-col gap-3">
        <p class="text-xs font-medium uppercase tracking-[0.2em] text-bloom">Flower</p>
        <h1 class="font-heading text-4xl font-semibold text-ink">Plan work the way it actually ships.</h1>
        <p class="max-w-xl text-base leading-7 text-ink-700">
          Flower is a task management workflow tool in the spirit of Pivotal Tracker. Stories move through icebox,
          backlog, current iteration, and done.
        </p>
      </div>

      <section class="rounded-xl border border-paper-200 bg-white p-5" data-testid="api-status">
        <h2 class="text-sm font-medium uppercase tracking-wider text-ink-700">API status</h2>
        <Show when={!health.loading && health()} fallback={<p class="mt-2 text-sm text-ink-600">Checking the API…</p>}>
          {(view) => <ApiStatusMessage view={view()} />}
        </Show>
      </section>

      <ul class="grid gap-3 sm:grid-cols-2">
        <For each={PANELS}>
          {(panel) => (
            <li class={cn("rounded-xl border border-paper-200 bg-white p-4")}>
              <h3 class="text-sm font-semibold text-ink">{panel.title}</h3>
              <p class="mt-1 text-sm leading-6 text-ink-700">{panel.body}</p>
            </li>
          )}
        </For>
      </ul>
    </main>
  );
}
