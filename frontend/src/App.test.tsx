/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { cleanup, render, screen } from "@solidjs/testing-library";
import App from "./App";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
  globalThis.fetch = mock(async () => {
    return new Response(JSON.stringify({ status: "ok", version: "dev" }), {
      status: 200,
      headers: { "content-type": "application/json" },
    });
  }) as unknown as typeof fetch;
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

describe("App", () => {
  it("renders the sign up route", async () => {
    render(() => <App />);
    expect(await screen.findByRole("heading", { name: "Sign up" })).toBeTruthy();
    expect(screen.getByText("Email and a password. That’s enough to start.")).toBeTruthy();
  });
});
