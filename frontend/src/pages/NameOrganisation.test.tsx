/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { NameOrganisationPage } from "./NameOrganisation";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

function meBody(verified: boolean) {
  return {
    id: "u1",
    username: "maya",
    email: "maya@example.com",
    display_name: "maya",
    email_verified_at: verified ? "2026-08-21T00:00:00Z" : null,
    organisations: [],
    last_project: null,
  };
}

describe("NameOrganisationPage", () => {
  it("asks for the organisation name after /me", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify(meBody(true)), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    render(() => (
      <Router>
        <Route path="/" component={NameOrganisationPage} />
      </Router>
    ));
    expect(screen.queryByRole("button", { name: "Create organisation" })).toBeNull();
    expect(await screen.findByText("What’s the organisation called?")).toBeTruthy();
    expect(screen.getByLabelText("Organisation name")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Create organisation" })).toBeTruthy();
  });

  it("does not flash Create organisation for an unverified user", async () => {
    let resolveMe: (value: Response) => void = () => undefined;
    globalThis.fetch = mock(() => {
      return new Promise<Response>((resolve) => {
        resolveMe = resolve;
      });
    }) as unknown as typeof fetch;
    render(() => (
      <Router>
        <Route path="/" component={NameOrganisationPage} />
      </Router>
    ));
    expect(screen.queryByRole("button", { name: "Create organisation" })).toBeNull();
    resolveMe(
      new Response(JSON.stringify(meBody(false)), {
        status: 200,
        headers: { "content-type": "application/json" },
      }),
    );
    expect(await screen.findByTestId("verify-block")).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Create organisation" })).toBeNull();
  });

  it("resends verify-email, not a magic link", async () => {
    const calls: string[] = [];
    globalThis.fetch = mock(async (input: RequestInfo | URL) => {
      const url = String(input);
      calls.push(url);
      if (url.includes("/auth/verify-email") && !url.includes("consume")) {
        return new Response(JSON.stringify({ ok: true }), {
          status: 200,
          headers: { "content-type": "application/json" },
        });
      }
      return new Response(JSON.stringify(meBody(false)), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    render(() => (
      <Router>
        <Route path="/" component={NameOrganisationPage} />
      </Router>
    ));
    fireEvent.click(await screen.findByRole("button", { name: "Resend the email" }));
    expect(await screen.findByText("Check your email for a link.")).toBeTruthy();
    expect(calls.some((url) => url.includes("/auth/verify-email") && !url.includes("consume"))).toBe(true);
    expect(calls.some((url) => url.includes("/auth/magic-link"))).toBe(false);
  });
});
