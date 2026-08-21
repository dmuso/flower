/** @jsxImportSource solid-js */
import { afterEach, beforeEach, describe, expect, it, mock } from "bun:test";
import { Router, Route } from "@solidjs/router";
import { cleanup, render, screen } from "@solidjs/testing-library";
import { BoardPage } from "./Board";

const originalFetch = globalThis.fetch;

beforeEach(() => {
  cleanup();
  window.history.replaceState({}, "", "/projects/proj-1");
});

afterEach(() => {
  cleanup();
  globalThis.fetch = originalFetch;
});

function renderBoard() {
  return render(() => (
    <Router>
      <Route path="/projects/:projectId" component={BoardPage} />
      <Route path="/signin" component={() => <p>sign in page</p>} />
      <Route path="/organisations/new" component={() => <p>new org</p>} />
      <Route path="/organisations/:organisationId/projects/new" component={() => <p>new project</p>} />
    </Router>
  ));
}

function json(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "content-type": "application/json" },
  });
}

function meBody(partial: Record<string, unknown> = {}) {
  return {
    id: "u1",
    username: "maya",
    email: "maya@example.com",
    display_name: "maya",
    email_verified_at: "2026-08-21T00:00:00Z",
    organisations: [],
    last_project: null,
    ...partial,
  };
}

function mockMissingBoard(me: unknown) {
  globalThis.fetch = mock(async (input: RequestInfo | URL) => {
    const url = String(input);
    if (url.includes("/api/v1/me")) {
      return json(me);
    }
    return json({ error: { code: "not_found", message: "We can’t find that." } }, 404);
  }) as unknown as typeof fetch;
}

describe("BoardPage", () => {
  it("exports the board page", () => {
    expect(typeof BoardPage).toBe("function");
  });

  it("shows paper column skeletons on first paint, not Loading…", () => {
    globalThis.fetch = mock(() => new Promise(() => undefined)) as unknown as typeof fetch;
    renderBoard();
    expect(screen.queryByText("Loading…")).toBeNull();
    expect(screen.getByTestId("board-skeletons")).toBeTruthy();
    expect(screen.getByTestId("skeleton-icebox")).toBeTruthy();
    expect(screen.getByTestId("skeleton-backlog")).toBeTruthy();
    expect(screen.getByTestId("skeleton-current")).toBeTruthy();
    expect(screen.getByTestId("skeleton-done")).toBeTruthy();
  });

  it("shows signed-out copy on 401", async () => {
    globalThis.fetch = mock(async () => {
      return json({ error: { code: "unauthorized", message: "You’re signed out." } }, 401);
    }) as unknown as typeof fetch;
    renderBoard();
    expect(await screen.findByText("You’re signed out.")).toBeTruthy();
    expect(screen.getByRole("link", { name: "Sign in" })).toBeTruthy();
    expect(screen.queryByText("Loading…")).toBeNull();
  });

  it("shows missing copy on 404 and falls back to new organisation", async () => {
    globalThis.fetch = mock(async () => {
      return json({ error: { code: "not_found", message: "We can’t find that." } }, 404);
    }) as unknown as typeof fetch;
    renderBoard();
    expect(await screen.findByText("We can’t find that.")).toBeTruthy();
    const back = screen.getByRole("link", { name: "Back to your projects" });
    expect(back).toBeTruthy();
    expect(back.getAttribute("href")).toBe("/organisations/new");
    expect(screen.queryByText("Loading…")).toBeNull();
  });

  it("sends 404 back to the last project", async () => {
    mockMissingBoard(
      meBody({
        last_project: { id: "p1", organisation_id: "o1", name: "Trail", slug: "trail" },
      }),
    );
    renderBoard();
    const back = await screen.findByRole("link", { name: "Back to your projects" });
    expect(screen.getByText("We can’t find that.")).toBeTruthy();
    expect(back.getAttribute("href")).toBe("/projects/p1");
  });

  it("sends 404 back to name a project when an organisation exists", async () => {
    mockMissingBoard(meBody({ organisations: [{ id: "o1", name: "Acme", role: "owner" }] }));
    renderBoard();
    const back = await screen.findByRole("link", { name: "Back to your projects" });
    expect(back.getAttribute("href")).toBe("/organisations/o1/projects/new");
  });

  it("sends 404 back to name an organisation when none exist", async () => {
    mockMissingBoard(meBody());
    renderBoard();
    const back = await screen.findByRole("link", { name: "Back to your projects" });
    expect(back.getAttribute("href")).toBe("/organisations/new");
  });
});
