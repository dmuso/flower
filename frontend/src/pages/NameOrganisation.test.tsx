/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { NameOrganisationPage } from "./NameOrganisation";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

describe("NameOrganisationPage", () => {
  it("asks for the organisation name", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(
        JSON.stringify({
          id: "u1",
          username: "maya",
          email: "maya@example.com",
          display_name: "maya",
          email_verified_at: "2026-08-21T00:00:00Z",
          organisations: [],
          last_project: null,
        }),
        { status: 200, headers: { "content-type": "application/json" } },
      );
    }) as unknown as typeof fetch;
    render(() => (
      <Router>
        <Route path="/" component={NameOrganisationPage} />
      </Router>
    ));
    expect(await screen.findByText("What’s the organisation called?")).toBeTruthy();
    expect(screen.getByLabelText("Organisation name")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Create organisation" })).toBeTruthy();
  });
});
