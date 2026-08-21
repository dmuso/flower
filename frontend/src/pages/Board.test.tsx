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
    </Router>
  ));
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
      return new Response(JSON.stringify({ error: { code: "unauthorized", message: "You’re signed out." } }), {
        status: 401,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    renderBoard();
    expect(await screen.findByText("You’re signed out.")).toBeTruthy();
    expect(screen.getByRole("link", { name: "Sign in" })).toBeTruthy();
    expect(screen.queryByText("Loading…")).toBeNull();
  });

  it("shows missing copy on 404", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ error: { code: "not_found", message: "We can’t find that." } }), {
        status: 404,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    renderBoard();
    expect(await screen.findByText("We can’t find that.")).toBeTruthy();
    expect(screen.getByRole("link", { name: "Back to your projects" })).toBeTruthy();
    expect(screen.queryByText("Loading…")).toBeNull();
  });
});
