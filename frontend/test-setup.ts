import { afterEach, beforeEach } from "bun:test";
import { Window } from "happy-dom";
import { solidPlugin } from "./scripts/plugins/solid-plugin";

process.env.FRONTEND_API_URL = "http://localhost:8180";

const happyWindow = new Window({ url: "http://localhost/" });
// Assign global references early so Solid and testing-library pick up the DOM.
globalThis.window = happyWindow as unknown as Window & typeof globalThis;
globalThis.document = happyWindow.document as unknown as Document;
globalThis.navigator = happyWindow.navigator as unknown as Navigator;
globalThis.HTMLElement = happyWindow.HTMLElement as unknown as typeof HTMLElement;
globalThis.Element = happyWindow.Element as unknown as typeof Element;
globalThis.Node = happyWindow.Node as unknown as typeof Node;
globalThis.Event = happyWindow.Event as unknown as typeof Event;
if (typeof globalThis.getComputedStyle !== "function") {
  globalThis.getComputedStyle = happyWindow.getComputedStyle.bind(happyWindow) as typeof getComputedStyle;
}

if (!(globalThis as { __flowerSolidPluginRegistered?: boolean }).__flowerSolidPluginRegistered) {
  Bun.plugin(solidPlugin);
  (globalThis as { __flowerSolidPluginRegistered?: boolean }).__flowerSolidPluginRegistered = true;
}

const solidCore = await import("solid-js/dist/dev.js");

const { mock } = await import("bun:test");
mock.module("solid-js", () => solidCore);
mock.module("solid-js/dist/server.js", () => solidCore);
mock.module("solid-js/dist/server.cjs", () => solidCore);
mock.module("solid-js/dist/solid.js", () => solidCore);
mock.module("solid-js/dist/solid.cjs", () => solidCore);

const solidWeb = await import("solid-js/web/dist/dev.js");

mock.module("solid-js/web", () => ({
  ...solidWeb,
  Portal: (props: { children: unknown }) => props.children,
}));
mock.module("solid-js/web/dist/server.js", () => ({
  ...solidWeb,
  Portal: (props: { children: unknown }) => props.children,
}));
mock.module("solid-js/web/dist/server.cjs", () => ({
  ...solidWeb,
  Portal: (props: { children: unknown }) => props.children,
}));

if (!("ResizeObserver" in globalThis)) {
  globalThis.ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  } as unknown as typeof ResizeObserver;
}

if (!("requestAnimationFrame" in globalThis)) {
  let nextAnimationFrameHandle = 1;
  const activeAnimationFrames = new Set<number>();

  globalThis.requestAnimationFrame = ((callback: FrameRequestCallback) => {
    const handle = nextAnimationFrameHandle++;
    activeAnimationFrames.add(handle);
    queueMicrotask(() => {
      if (!activeAnimationFrames.delete(handle)) {
        return;
      }
      callback(Date.now());
    });
    return handle;
  }) as typeof requestAnimationFrame;

  globalThis.cancelAnimationFrame = ((handle: number) => {
    activeAnimationFrames.delete(handle);
  }) as typeof cancelAnimationFrame;
}

const { cleanup } = await import("@solidjs/testing-library");

const blockedFetch = (async (input: RequestInfo | URL) => {
  throw new Error(
    `Unexpected live network request in frontend test: ${String(input)}. Mock the API boundary instead of calling a real service.`,
  );
}) as typeof fetch;

beforeEach(() => {
  globalThis.fetch = blockedFetch;
});

afterEach(() => {
  cleanup();
});
