/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { HomePage } from "./Home";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

describe("HomePage", () => {
  it("renders the product name and planning panels", () => {
    globalThis.fetch = mock(async () => {
      return new Promise<Response>(() => undefined);
    }) as unknown as typeof fetch;

    render(() => <HomePage />);

    expect(screen.getByRole("heading", { name: "Plan work the way it actually ships." })).toBeTruthy();
    expect(screen.getByText("Icebox")).toBeTruthy();
    expect(screen.getByText("Backlog")).toBeTruthy();
    expect(screen.getByText("Current")).toBeTruthy();
    expect(screen.getByText("Done")).toBeTruthy();
    expect(screen.getByTestId("api-status")).toBeTruthy();
  });

  it("shows a connected status when the API is healthy", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ status: "ok", version: "dev" }), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;

    render(() => <HomePage />);

    const status = await screen.findByTestId("api-status-ok");
    expect(status.textContent).toContain("Connected to Flower API dev");
  });

  it("shows an error when the API is unreachable", async () => {
    globalThis.fetch = mock(async () => {
      return new Response("unavailable", { status: 503 });
    }) as unknown as typeof fetch;

    render(() => <HomePage />);

    const status = await screen.findByTestId("api-status-error");
    expect(status.textContent).toContain("The API is not reachable.");
  });
});
